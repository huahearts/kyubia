package main

import "github.com/huahearts/kyubia/knet"

func main() {
	s := knet.NewServer()
	s.Start()
}
