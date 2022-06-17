package kiface
type IRequest struct {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
}