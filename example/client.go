package main

import (
	"io"
	"net"

	"github.com/shunhui19/goes/lib"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		lib.Fatal("connect server error: ", err)
	}
	defer conn.Close()

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil || err == io.EOF {
				lib.Fatal("client read error: %v", err)
			}
			lib.Info("Received content: %v", string(buf[:n]))
		}
	}()

	content := "hello, world\n"
	_, err = conn.Write([]byte(content))
	if err != nil {
		lib.Warn("client send error: %v", err)
	}
	lib.Info("client send content: %v", content)

	select {}
}
