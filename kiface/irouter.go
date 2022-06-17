package kiface

type IRouter interface {
	PreCallback(req IRequest)
	Callback(req IRequest)
	PostCallback(req IRequest)
}
