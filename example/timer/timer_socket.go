package main

import (
	"time"

	"github.com/shunhui19/goes"
	"github.com/shunhui19/goes/connections"
	"github.com/shunhui19/goes/lib"
	"github.com/shunhui19/goes/protocols"
)

func main() {
	g := goes.NewGoer("127.0.0.1:9090", protocols.NewTextProtocol(), "tcp")
	t := lib.NewTimer(60, 1*time.Second)

	g.OnMessage = func(connection connections.Connection, data []byte) {
		// 每2秒给客户端发送消息
		t.Add(2*time.Second, func(v ...interface{}) {
			lib.Info("[%v]: send message: hello, world", time.Now().Format("2006-01-02 15:04:05.9999"))
			connection.Send("hello, world\n", false)
		}, connection, true)
	}

	// 启动定时器
	go t.Start()

	g.RunAll()
}
