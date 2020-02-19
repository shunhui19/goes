package goes

import (
	"bytes"
	"fmt"
	"goes/connections"
	"goes/protocols"
	"net"
	"testing"
)

var goes *Goer

func init() {
	goes = NewGoer("127.0.0.1:8080", protocols.NewTextProtocol(), "tcp")
	goes.OnConnect = func(connection connections.Connection) {
		fmt.Println("new client is coming")
	}
	goes.OnMessage = func(connection connections.Connection, data []byte) {
		connection.Send("Request received: "+string(data), false)
	}
	goes.OnClose = func() {
		fmt.Println("client is closed.")
	}
	go func() {
		goes.RunAll()
	}()
}

func TestGoer_RunAll(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Errorf("could not connect goer: %v", err)
	}
	if conn == nil {

	}
	defer conn.Close()
}

func TestGoer_OnCallback(t *testing.T) {
	tt := []struct {
		name    string
		message []byte
		want    []byte
	}{
		{
			"OnMessage callback one",
			[]byte("OnConnect one\n"),
			[]byte("Request received: OnConnect one\n"),
		},
		{
			"OnMessage callback two",
			[]byte("OnConnect\n"),
			[]byte("Request received: OnConnect\n"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				t.Errorf("connect goer fail: %v", err)
			}
			defer conn.Close()

			_, err = conn.Write(tc.message)
			if err != nil {
				t.Errorf("could not write message to goer: %v", err)
			}

			out := make([]byte, 1024)
			if n, err := conn.Read(out); err == nil {
				if bytes.Compare(out[:n], tc.want) != 0 {
					t.Error("response did not match expected output")
				}
			} else {
				t.Errorf("could not read message from connection.")
			}
		})
	}
}
