package knet

import "github.com/huahearts/kyubigo/kiface"

type Option func(s *Server)

func WithPacket(pack kiface.IPacket) Option {
	return func(s *Server) {
		s.packet = pack
	}
}
