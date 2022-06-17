package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type GOBJ struct {
	TCPServer        kiface.IServer
	Host             string
	TCPPort          int
	Name             string
	Version          string
	MaxPacketSize    int
	MaxConn          uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32

	ConfFilePath string

	LogDir        string
	LogFile       string
	LogDebugClose bool
}

var GlobalObject *GOBJ

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, nil
	}

	return false, err
}

func (g *GOBJ) Reload() {
	if confFileExists, _ := PathExist(g.ConfFilePath); confFileExists == false {
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

	if g.LogFile != "" {
		//klog.SetLogFile(g.LogDir,g.LogFile)
	}

	if g.LogDebugClose {
		//klog.CloseBug()
	}
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	GlobalObject := &GOBJ{
		Name:             "KyubiServer",
		Version:          "v1.0.0",
		TCPPort:          9981,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     pwd + "/conf/kyubi.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "",
		LogDebugClose:    false,
	}

	GlobalObject.Reload()
}
