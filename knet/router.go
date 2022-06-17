package knet

import "github.com/huahearts/kyubigo/kiface"

type BaseRouter struct{}

func (br *BaseRouter) PreCallback(req kiface.IRequest)  {}
func (br *BaseRouter) Callback(req kiface.IRequest)     {}
func (br *BaseRouter) PostCallback(req kiface.IRequest) {}
