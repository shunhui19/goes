package goes

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/shunhui19/goes/connections"
	"github.com/shunhui19/goes/protocols"
)

var tcp, tcp4, udp *Goer

func init() {
	tcp = NewGoer("127.0.0.1:8080", protocols.NewTextProtocol(), "tcp")
	tcp.OnConnect = func(connection connections.Connection) {
		fmt.Println("tcp client is coming")
	}
	tcp.OnMessage = func(connection connections.Connection, data []byte) {
		connection.Send("Request received: "+string(data), false)
	}
	tcp.OnClose = func(connection connections.Connection) {
		fmt.Println("tcp client is closed.")
	}

	tcp4 = NewGoer("127.0.0.1:8081", protocols.NewTextProtocol(), "tcp4")
	tcp4.OnConnect = func(connection connections.Connection) {
		fmt.Println("tcp4 client is coming")
	}
	tcp4.OnMessage = func(connection connections.Connection, data []byte) {
		connection.Send("Request received: "+string(data), false)
	}
	tcp4.OnClose = func(connection connections.Connection) {
		fmt.Println("tcp4 client is closed.")
	}

	udp := NewGoer("127.0.0.1:9090", protocols.NewTextProtocol(), "udp")
	udp.OnMessage = func(connection connections.Connection, data []byte) {
		connection.Send("Request received: "+string(data), false)
	}

	go tcp.RunAll()
	go tcp4.RunAll()
	go udp.RunAll()
}

func TestGoer_RunAll(t *testing.T) {
	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", ":8080"},
		{"tcp4", ":8081"},
		{"udp", ":9090"},
	}

	for _, server := range servers {
		conn, err := net.Dial(server.protocol, server.addr)
		if err != nil {
			t.Errorf("could not connect goer: %v", err)
		}
		if conn == nil {

		}
		defer conn.Close()
	}

}

func TestGoer_OnCallback(t *testing.T) {
	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", ":8080"},
		{"tcp4", ":8081"},
		{"udp", ":9090"},
	}

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

	for _, server := range servers {
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				conn, err := net.Dial(server.protocol, server.addr)
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
}

func BenchmarkGenerateConnectionID(t *testing.B) {
	for i := 0; i < t.N; i++ {
		tcp.generateConnectionID()
	}
}
