package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	var w sync.WaitGroup
	clientCount := 50
	sendCount := 5
	for i := 0; i < clientCount; i++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			// receive
			//go func() {
			//	buf := make([]byte, 1024)
			//	for {
			//		n, err := conn.Read(buf)
			//		if err != nil {
			//			log.Fatal(err.Error())
			//		}
			//		//fmt.Println(bytes.Contains(buf[:n], []byte("\n")))
			//		fmt.Println(string(buf[:n]))
			//	}
			//}()

			rand.Seed(time.Now().Unix())
			for j := 0; j < sendCount; j++ {
				n, err := conn.Write([]byte(fmt.Sprintf("[goroutine-%d]", i) + "hello" + strconv.Itoa(j) + "\n"))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("success send %v byte\n", n)
				//time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
			}

			// 保持客户端收到服务端数据之前不关闭连接
			//time.Sleep(5 * time.Second)
		}(i)
	}

	w.Wait()
}
