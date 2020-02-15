// connection interface.
package connections

// ConnectionInterface the method of interface.
type ConnectionInterface interface {
	// Send sends data on the connection.
	Send(data string, raw bool) interface{}
	// Close close connection.
	Close(data string)
	// GetRemoteIP get remote IP.
	GetRemoteIP() string
	// GetRemotePort get remote port.
	GetRemotePort() int
	// GetRemoteAddress get remote address.
	GetRemoteAddress() string
	// GetLocalIP get local IP.
	GetLocalIP() string
	// GetLocalPort get local port.
	GetLocalPort() int
	// GetLocalAddress get local address.
	GetLocalAddress() string
}
