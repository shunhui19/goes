package main

import (
	"goes"
	"goes/connections"
	"goes/lib"
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
	//goer.OnConnect = func(connection *connections.TCPConnection) {
	//	fmt.Printf("remoteAddress: %s\n", connection.GetRemoteAddress())
	//	fmt.Printf("remoteIp: %s\n", connection.GetRemoteIP())
	//	fmt.Printf("remotePort: %d\n", connection.GetRemotePort())
	//
	//	fmt.Printf("localAddress: %s\n", connection.GetLocalAddress())
	//	fmt.Printf("localIp: %s\n", connection.GetLocalIP())
	//	fmt.Printf("localPort: %d\n", connection.GetLocalPort())
	//}

	goer.OnConnect = func(connection connections.ConnectionInterface) {
		//fmt.Println("new connection is coming")
		//connection.Send("hello, world", false)
	}

	goer.OnMessage = func(connection connections.ConnectionInterface, data []byte) {
		lib.Info("client[%v], content: %v", connection.GetRemoteAddress(), string(data))
		connection.Send("ok", false)

		// send data to other client.
		goer.Connections.Range(func(key, value interface{}) bool {
			if c := value.(connections.ConnectionInterface); c.GetRemoteAddress() != connection.GetRemoteAddress() {
				c.Send("a new client is coming", false)
			}
			return true
		})
	}

	goer.OnClose = func() {
		//lib.Info("client is closed")
	}

	//goer.Transport = "udp"
	goer.RunAll()
}
