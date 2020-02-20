package main

import (
	"fmt"
	"goes"
	"goes/connections"
	"goes/lib"
	"goes/protocols"
)

func main() {
	goer := goes.NewGoer("127.0.0.1:8080", protocols.NewTextProtocol(), "tcp")
	goer.StdoutFile = "./goes.log"

	goer.OnMessage = func(connection connections.Connection, data []byte) {
		lib.Info("Receive request: %v", string(data))
		connection.Send("hello, goes", false)

		// notice other clients a new client is connect.
		goer.Connections.Range(func(key, value interface{}) bool {
			if c := value.(connections.Connection); c.GetRemoteAddress() != connection.GetRemoteAddress() {
				c.Send(fmt.Sprintf("[%v] client is connect", connection.GetLocalAddress()), false)
			}
			return true
		})
	}

	goer.OnClose = func(connection connections.Connection) {
		lib.Info("client is close")
	}

	goer.RunAll()
}
