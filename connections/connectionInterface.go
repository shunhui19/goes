// connection interface.
package connections

// Connection the method of interface.
type Connection interface {
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

// CStore the interface of store all connection.
type CStore interface {
	// Set store a TCPConnection.
	Set(conn *TCPConnection)
	// Get return a TCPConnection.
	Get(connID int) (*TCPConnection, bool)
	// Del remove a TCPConnection from sync.Map.
	Del(connID int)
	// Range calls f sequentially for each key and value present in the map.
	Range(f func(key, value interface{}) bool)
	// Len return the count of all TCPConnection.
	Len() int32
}
