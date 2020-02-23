package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/shunhui19/goes/lib"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var w sync.WaitGroup
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				lib.Fatal("connect server error: ", err)
			}
			defer conn.Close()

			w.Add(1)
			go receive(conn, &w)

			// 每个连接送100条数据
			for j := 0; j < 100; j++ {
				send(conn, []byte(fmt.Sprintf("[%d]%d-hello\n", i, j)))
			}

			w.Wait()
		}(i)
	}
	wg.Wait()
	lib.Info("done")
}

func receive(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, 1024)
	receivedCount := 0
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				lib.Warn("server close.")
				return
			}
			lib.Warn("client received error: %v\n", err)
			break
		}
		receivedCount++
		lib.Info("The %d received, content: %v", receivedCount, string(buf[:n]))
	}
}

func send(conn net.Conn, data []byte) {
	n, err := conn.Write(data)
	if err != nil {
		lib.Warn("send error: %v", err)
	}
	lib.Info("Send %v bytes\n", n)
}
