package main

import "github.com/huahearts/kyubigo/knet"

func main() {
	s := knet.NewServer()
	s.Start()
}
