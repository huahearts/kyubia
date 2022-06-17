package knet

import "github.com/huahearts/kyubia/kiface"

type Option func(s *Server)

func WithPacket(pack kiface.IPacket) Option {
	return func(s *Server) {
		s.packet = pack
	}
}
