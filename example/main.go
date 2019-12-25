package main

import (
	"fmt"
	"goes"
	"goes/connections"
	"goes/protocols"
)

func main() {
	goer := goes.NewGoer("127.0.0.1:8080", protocols.NewTextProtocol(), "udp")
	//goer.OnConnect = func(connection *connections.TcpConnection) {
	//	fmt.Printf("remoteAddress: %s\n", connection.GetRemoteAddress())
	//	fmt.Printf("remoteIp: %s\n", connection.GetRemoteIp())
	//	fmt.Printf("remotePort: %d\n", connection.GetRemotePort())
	//
	//	fmt.Printf("localAddress: %s\n", connection.GetLocalAddress())
	//	fmt.Printf("localIp: %s\n", connection.GetLocalIp())
	//	fmt.Printf("localPort: %d\n", connection.GetLocalPort())
	//}

	goer.OnConnect = func(connection connections.ConnectionInterface) {
		connection.Send("hello, world", false)
	}

	goer.OnMessage = func(connection connections.ConnectionInterface, data []byte) {
		fmt.Println(string(data))
		connection.Send(fmt.Sprintf("the client send message is %v", string(data)), false)
		//fmt.Println(status)
	}

	goer.OnClose = func() {
		fmt.Println("closed")
	}

	//goer.Transport = "udp"
	goer.RunAll()
}
