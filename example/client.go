package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 5; i++ {
		n, err := conn.Write([]byte("hello" + strconv.Itoa(i)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("success send %v byte\n", n)
		time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
	}
}
