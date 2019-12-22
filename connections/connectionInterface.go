// connection interface.
package connections

type ConnectionInterface interface {
	// Send sends data on the connection.
	Send(data string, raw bool) interface{}
	// Close close connection.
	Close(data string)
	// GetRemoteIp get remote IP.
	GetRemoteIp() string
	// GetRemotePort get remote port.
	GetRemotePort() int
	// GetRemoteAddress get remote address.
	GetRemoteAddress() string
	// GetLocalIp get local IP.
	GetLocalIp() string
	// GetLocalPort get local port.
	GetLocalPort() int
	// GetLocalAddress get local address.
	GetLocalAddress() string
}
