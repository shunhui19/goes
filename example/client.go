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
	for i := 0; i < 500; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal(err)
			}

			// receive
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						log.Fatal(err.Error())
					}
					//fmt.Println(bytes.Contains(buf[:n], []byte("\n")))
					fmt.Println(string(buf[:n]))
				}
			}()

			rand.Seed(time.Now().Unix())
			for i := 0; i < 10; i++ {
				n, err := conn.Write([]byte("hello" + strconv.Itoa(i) + "\n"))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("success send %v byte\n", n)
				//time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
			}

			//time.Sleep(10 * time.Second)
		}()
	}

	w.Wait()
}
