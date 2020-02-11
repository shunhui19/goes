package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var w sync.WaitGroup
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal("connect server error: ", err)
			}
			defer conn.Close()

			w.Add(1)
			go receive(conn, &w)

			// 每个连接送数据
			for j := 0; j < 10000; j++ {
				send(conn, []byte(fmt.Sprintf("[%d]%d-hello\n", i, j)))
				//time.Sleep(time.Second)
			}

			w.Wait()
		}(i)
	}
	wg.Wait()
	fmt.Println("done")
}

func receive(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, 1024)
	receivedCount := 0
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Fatalln("server close.")
			}
			log.Printf("client received error: %v\n", err)
			continue
		}
		receivedCount++
		//fmt.Printf("The %d received, content: %v", receivedCount, string(buf[:n]))
	}
}

func send(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Send %v bytes\n", n)
}
