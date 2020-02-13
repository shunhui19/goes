package http

import (
	"testing"
)

func TestHttp_getRequestSize(t *testing.T) {
	type args struct {
		header []byte
		method string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "GET", args: args{header: []byte("GET / HTTP/1.1\r\nHost: www.frank.com\r\nConnection: keep-live\r\n"), method: "GET"}, want: 64},
		{name: "OPTIONS", args: args{header: []byte("OPTIONS / HTTP/1.1\r\nHost: www.frank.com\r\nConnection: keep-live\r\n"), method: "OPTIONS"}, want: 68},
		{name: "HEAD", args: args{header: []byte("HEAD / HTTP/1.1\r\nHost: www.frank.com\r\nConnection: keep-live\r\n"), method: "HEAD"}, want: 65},
		{name: "DELETE", args: args{header: []byte("DELETE / HTTP/1.1\r\nHost: www.frank.com\r\nConnection: keep-live\r\n"), method: "DELETE"}, want: 67},
		{name: "POST", args: args{header: []byte("POST / HTTP/1.1\r\nHost: www.frank.com\r\nContent-Length: 15\r\nConnection: keep-live\r\n"), method: "POST"}, want: 19},
		{name: "PUT", args: args{header: []byte("PUT / HTTP/1.1\r\nHost: www.frank.com\r\nContent-Length: 15\r\nConnection: keep-live\r\n"), method: "PUT"}, want: 19},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Http{}
			if got := h.getRequestSize(tt.args.header, tt.args.method); got != tt.want {
				t.Errorf("getRequestSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttp_Decode(t *testing.T) {
	type args struct {
		packet []byte
	}
	tests := []struct {
		args
	}{
		{args{packet: []byte("GET /helloWorld HTTP/1.1\r\nHOST:www.frank.com\r\nCookie: BAIDUID=C88B6F1243FCB8D658349758E3340CB8:FG=1; BIDUPSID=C88B6F1243FCB8D658349758E3340CB8; PSTM=1556635667\r\n\r\n")}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			h := &Http{}
			h.Decode(tt.args.packet)
		})
	}
}

func TestHttp_Encode(t *testing.T) {
	type args struct {
		h Http
	}
	tests := []struct {
		args
	}{
		{Http{header: map[string]interface{}{"Content-Type": "text/html"}, gzip: false}},
	}
}
