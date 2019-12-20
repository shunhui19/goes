package main

import (
	"fmt"
	"goes/lib"
	"net"
)

var buildInProtocol = []map[string]string{
	{"tcp": "tcp"},
	{"tcp4": "tcp4"},
	{"tcp6": "tcp6"},
	{"udp": "udp"},
	{"unix": "unix"},
	{"unixpacket": "unixpacket"},
	{"ssl": "tcp"},
}

type Gs struct {
	protocol   string
	socketName string
	mainSocket net.Listener
}

func NewGs(protocol string, address string) *Gs {
	if protocol == "" || address == "" {
		lib.Fatal("protocol:%s and address:%s can not bu null", protocol, address)
	}

	exist := false
	for _, v := range buildInProtocol {
		if v[protocol] == protocol {
			exist = true
			break
		}
	}
	if !exist {
		lib.Fatal("unsupported protocol: %v", protocol)
	}

	listener, err := net.Listen(protocol, address)
	if err != nil {
		lib.Fatal(err.Error())
	}
	lib.Info("server start...")

	return &Gs{
		protocol:   protocol,
		socketName: address,
		mainSocket: listener,
	}
}

func main() {
	//listener, err := net.Listen("tcp4", "127.0.0.1:8080")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("server start...")
	goer := NewGs("tcp4", "127.0.0.1:8080")

	b := make([]byte, 100)
	for {
		conn, err := goer.mainSocket.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func() {
			defer conn.Close()
			for {
				n, _ := conn.Read(b)
				fmt.Println(b)
				// 没有读到数据，表示客户端断开了连接
				if n == 0 {
					fmt.Println("client is closed")
					break
				} else {
					fmt.Println("a new client is coming")
					fmt.Printf("content: %v\n", string(b[:n]))
				}
			}
		}()
	}
}
