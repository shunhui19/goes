package main

import (
	"fmt"
	"goes"
	"goes/connections"
)

func main() {
	goer := goes.NewGoer("127.0.0.1:8080", nil, "")
	goer.OnConnect = func(connection *connections.TcpConnection) {
		fmt.Printf("remoteAddress: %s\n", connection.GetRemoteAddress())
		fmt.Printf("remoteIp: %s\n", connection.GetRemoteIp())
		fmt.Printf("remotePort: %d\n", connection.GetRemotePort())

		fmt.Printf("localAddress: %s\n", connection.GetLocalAddress())
		fmt.Printf("localIp: %s\n", connection.GetLocalIp())
		fmt.Printf("localPort: %d\n", connection.GetLocalPort())
	}

	goer.OnMessage = func(connection *connections.TcpConnection, data string) {
		fmt.Println(data, "1111")
	}

	goer.OnClose = func() {
		fmt.Println("closed")
	}

	//goer.Transport = "udp"
	goer.RunAll()
}
