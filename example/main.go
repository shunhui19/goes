package main

import (
	"goes"
	"goes/connections"
	"goes/protocols"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("./cpuprofile.prof")
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
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
		//fmt.Printf("mainGoroutine[%d]Receive: %s\n", os.Getpid(), string(data))
		//lib.Info(string(data))
		connection.Send("HTTP/1.1 200 OK\r\nConnection: keep-alive\r\nServer: workerman\\1.1.4\r\n\r\nhello", false)
		//connection.Send(fmt.Sprintf("the client send message is %v", string(data)), false)
		//fmt.Println(status)
	}

	goer.OnClose = func() {
		//lib.Info("client is closed")
	}

	//goer.Transport = "udp"
	goer.RunAll()
}
