package knet

import (
	"fmt"
	"net"

	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/utils"
)

var zinxLogo = `                                        
	Kyubi
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        uint32
	msgHandler  kiface.IMsgHandler
	ConnMgr     kiface.IConnMgr
	OnConnStart func(conn kiface.IConnection)
	OnConnStop  func(conn kiface.IConnection)
	packet      kiface.IPacket
}

func NewServer(opts ...Option) kiface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       uint32(utils.GlobalObject.TCPPort),
		msgHandler: NewMsgHandler(),
		ConnMgr:    NewConnMgr(),
		packet:     NewDataPacket(),
	}

	/*for opt := range opts {
		opt(s)
	}*/
	return s
}

func (s *Server) Start() {
	printLogo()
	fmt.Printf("[Server Start] Server Name:%v,IPVersion:%v, IP:%v,Port:%v\n", s.Name, s.IPVersion, s.IP, s.Port)

	go func() {
		//工作池后续添加
		s.msgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen error")
			panic(err)
		}

		fmt.Println("start kyubi server ", s.Name, "succ,now listenning...")
		// 自动生成ID的方法
		var cID uint32
		cID = 0
		for { //while(1)循环
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			if uint32(s.ConnMgr.Len()) > utils.GlobalObject.MaxConn {
				//超出最大连接数 应该给客户端返回消息
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cID, s.msgHandler)
			cID++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP]kyubi server,name", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgId uint32, router kiface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() kiface.IConnMgr {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(onStartCallback func(kiface.IConnection)) {
	s.OnConnStart = onStartCallback
}
func (s *Server) SetOnConnStop(onStopCallback func(kiface.IConnection)) {
	s.OnConnStop = onStopCallback
}
func (s *Server) OnConnStartCallback(conn kiface.IConnection) {
	if s.OnConnStartCallback != nil {
		fmt.Println("--->callonStart callback")
		s.OnConnStartCallback(conn)
	}
}
func (s *Server) OnConnStopCallback(conn kiface.IConnection) {
	if s.OnConnStopCallback != nil {
		fmt.Println("--->callonStop callback")
		s.OnConnStartCallback(conn)
	}
}
func (s *Server) Packet() kiface.IPacket {
	return s.packet
}

func (s *Server) display() {
	fmt.Println("")
}
func printLogo() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/huahearts                    %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Kyubi] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
}
func init() {}
