package main

import (
	"fmt"
	"goes"
	"goes/connections"
	"goes/lib"
	"goes/protocols"
	"os"
)

func main() {
	goer := goes.NewGoer("127.0.0.1:8080", protocols.NewTextProtocol(), "tcp")
	goer.StdoutFile = "./server.log"
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
		//fmt.Println("new connection is coming")
		//connection.Send("hello, world", false)
	}

	goer.OnMessage = func(connection connections.ConnectionInterface, data []byte) {
		fmt.Printf("mainGoroutine[%d]Receive: %s\n", os.Getpid(), string(data))
		connection.Send(string(data), false)
		//connection.Send(fmt.Sprintf("the client send message is %v", string(data)), false)
		//fmt.Println(status)
	}

	goer.OnClose = func() {
		lib.Info("client is closed")
	}

	//goer.Transport = "udp"
	goer.RunAll()
}
