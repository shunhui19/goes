package main

import (
	"fmt"
	"goes/lib"
	"net"
	"os"
	"time"
)

func main() {
	done := make(chan struct{})
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		lib.Fatal("connect fail: %v", err)
	}
	defer conn.Close()

	// received message
	go func(conn net.Conn) {
		for {
			select {
			default:
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					lib.Warn("read content error: %v", err)
					return
				}
				fmt.Println("Receive:", string(buf[:n]))
			case <-done:
				return
			}
		}
	}(conn)

	_, err = conn.Write([]byte("client message\n"))
	if err != nil {
		lib.Warn("client send data error:%v", err)
		os.Exit(1)
	}

	time.Sleep(100 * time.Second)
	done <- struct{}{}
}
