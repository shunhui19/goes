package goes

import (
	"fmt"
	"goes/connections"
	"goes/protocols"
	"sync"
	"testing"
)

func TestNewGoer(t *testing.T) {
	type args struct {
		socketName          string
		applicationProtocol protocols.Protocol
		transportProtocol   string
		daemon              bool
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "goer-tcp", args: args{"127.0.0.1:8081", nil, "tcp", true}},
		{name: "goer-tcp4", args: args{"127.0.0.1:8082", nil, "tcp4", false}},
		//{name: "goer-tcp6", args: args{"127.0.0.1:8083", nil, "tcp6"}},
	}
	var w sync.WaitGroup
	for _, tt := range tests {
		w.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoer(tt.args.socketName, tt.args.applicationProtocol, tt.args.transportProtocol)
			go func() {
				defer w.Done()
				g.RunAll()
				g.OnConnect = func(connection connections.ConnectionInterface) {
					fmt.Println("OnConnect callback.")
				}
				g.OnMessage = func(connection connections.ConnectionInterface, data []byte) {
					fmt.Printf("OnMessage: %s\n", data)
				}
				g.OnClose = func() {
					fmt.Println("OnClose")
				}
			}()
		})
	}
	w.Wait()
}

// test start model, debug mode and daemon mode.
//func TestGoer_StartModel(t *testing.T) {
//	type args struct {
//		daemon bool
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{name: "daemon-model", args: args{true}},
//		{name: "debug-model", args: args{false}},
//		//{name: "goer-tcp6", args: args{"127.0.0.1:8083", nil, "tcp6"}},
//	}
//
//	for tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//		})
//	}
//	goer := NewGoer("127.0.0.1:9090", protocols.NewTextProtocol(), "tcp")
//	goer.Daemon = true
//}
