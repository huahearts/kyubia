package knet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/utils"
)

type Connection struct {
	TCPServer   kiface.IServer
	Conn        *net.TCPConn
	ConnId      uint32
	MsgHandler  kiface.IMsgHandler
	ctx         context.Context
	cancel      context.CancelFunc
	msgChan     chan []byte // 无缓冲管道 读写两个goroutine之间消息通信
	msgBuffChan chan []byte // 有缓冲
	sync.RWMutex

	property     map[string]interface{}
	propertyLock sync.Mutex

	isClosed bool
}

func NewConnection(server kiface.IServer, conn *net.TCPConn, connId uint32, msgHandler kiface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnId:      connId,
		MsgHandler:  msgHandler,
		isClosed:    false,
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    nil,
	}

	c.TCPServer.GetConnMgr().Add(c)
	return c
}

//写goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[writer goroutine is running]")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "conn writer exit")
	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send buff data errpr:", err, "Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//读goroutine
func (c *Connection) StartReader() {
	fmt.Println("[Reader goroutine is running]")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "conn reader exit")
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//read pack header
			headData := make([]byte, c.TCPServer.Packet().GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}

			msg, err := c.TCPServer.Packet().Unpack(headData)
			if err != nil {
				fmt.Println("unpack error", err)
				return
			}

			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read data err", err)
					return
				}
			}

			msg.SetData(data)

			req := &Request{
				conn: c,
				msg:  msg,
			}

			if utils.GlobalObject.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(req)
			} else {
				go c.MsgHandler.DoMsgHandler(req)
			}
		}
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.StartReader()

	go c.StartWriter()

	c.TCPServer.OnConnStartCallback(c)

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

func (c *Connection) Stop() {
	c.cancel()
}

func (c *Connection) GetTCPCOnnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnId
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接发送消息
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		return errors.New("send msg packet error")
	}

	_, err = c.Conn.Write(msg)
	return err
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()
	if c.isClosed {
		errors.New("Connection closed when send buff msg")
	}

	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(NewMsgPacket(msgID, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgID)
		return errors.New("pack error msg")
	}

	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) Context() *context.Context {
	return &c.ctx
}

func (c *Connection) finalizer() {
	c.TCPServer.OnConnStopCallback(c)
	c.Lock()
	defer c.Unlock()

	if c.isClosed {
		return
	}
	fmt.Println("Conn Stop()...ConnID = ", c.GetConnID())

	_ = c.Conn.Close()

	c.TCPServer.GetConnMgr().Remove(c)
	close(c.msgBuffChan)
	c.isClosed = true
}
