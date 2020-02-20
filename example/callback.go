package main

import (
	"goes"
	"goes/connections"
	"goes/lib"
)

func main() {
	goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")
	// 当有客户端与服务端完成三次握手之后
	goer.OnConnect = func(connection connections.Connection) {
		lib.Info("a new client is connected, addr: %v", connection.GetRemoteAddress())
	}
	// 当服务端接收客户端发送的消息
	goer.OnMessage = func(connection connections.Connection, data []byte) {
		lib.Info("client send message: %v", string(data))
	}
	// 应用层发送缓冲区满时
	goer.OnBufferFull = func(connection connections.Connection) {
		lib.Info("send buff is full")
	}
	// 应用层发送缓冲区为空时
	goer.OnBufferDrain = func(connection connections.Connection) {
		lib.Info("the send buf is going drain")
	}
	// 客户端连接断开时
	goer.OnClose = func(connection connections.Connection) {
		lib.Info("client is closed")
	}
	// 服务端准备关闭服务时
	goer.OnStop = func() {
		lib.Info("goer server is stopped")
	}
	goer.RunAll()
}
