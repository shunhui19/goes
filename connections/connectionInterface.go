// connection interface.
package connections

type connectionInterface interface {
	// Send sends data on the connection.
	Send(data string)
	// Close close connection.
	Close()
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
