// TcpConnection.
package connections

import (
	"net"
	"sync"
)

const (
	// READ_BUFFER_SEIZE read buffer size.
	READ_BUFFER_SEIZE = 65535
	// MAX_SEND_BUFFER_SIZE maximum size of send buffer.
	MAX_SEND_BUFFER_SIZE = 104856
	// MAX_SEND_BUFFER_SIZE default maximum size of send buffer.
	DEFAULT_MAX_SEND_BUFFER_SIZE = 104856
	// MAX_PACKAGE_SIZE maximum size of acceptable buffer.
	MAX_PACKAGE_SIZE = 104856
	// MAX_PACKAGE_SIZE default maximum size of acceptable buffer.
	DEFAULT_MAX_PACKAGE_SIZE = 1048560
)

// TcpConnection struct.
type TcpConnection struct {
	// baseConnection.
	baseConnection BaseConnection
	// OnMessage emitted when data is received.
	OnMessage func()
	// OnError emitted when a error occurs with connection.
	OnError func()
	// OnClose emitted when the other end of the socket send a FIN package.
	OnClose func()
	// OnBuffFull emitted when the send buffer becomes full.
	OnBuffFull func()
	// OnBufferDrain emitted when the send buffer becomes empty.
	OnBufferDrain func()
	// Protocol application layer protocol.
	Protocol string
	// Id the id of connection.
	Id int
	// MaxSendBufferSize set the maximum send buffer size for the current connection.
	MaxSendBufferSize int
	// MaxPackageSize set the maximum acceptable packet size for the current connection.
	MaxPackageSize int
	// socket tcp socket.
	socket *net.Conn
	// remoteAddress remote address.
	remoteAddress string
	// recvBuffer receive buffer.
	recvBuffer string
	// connections all connection instances, key is connection id and value is *net.Conn.
	connections sync.Map
}

// Send send data on the connection.
func (t *TcpConnection) Send(data string) {

}

// Close close connection.
func (t *TcpConnection) Close() {

}

// GetRemoteIp get remote ip.
func (t *TcpConnection) GetRemoteIp() string {
	return ""
}

// GetRemotePort get remote port.
func (t *TcpConnection) GetRemotePort() int {
	return 0
}

// GetRemoteAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *TcpConnection) GetRemoteAddress() string {
	return ""
}

// GetLocalIp get local ip.
func (t *TcpConnection) GetLocalIp() string {
	return ""
}

// GetLocalPort get local port.
func (t *TcpConnection) GetLocalPort() int {
	return 0
}

// GetLocalAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *TcpConnection) GetLocalAddress() string {
	return ""
}
