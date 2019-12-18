// updConnection.
package connections

import "net"

// UpdConnection struct.
type UdpConnection struct {
	// Protocol Application layer protocol.
	Protocol string
	// socket udp socket.
	socket net.PacketConn
	// remoteAddress remote address.
	remoteAddress string
}

// Send send data on the connection.
func (t *UdpConnection) Send(data string) {

}

// Close close connection.
func (t *UdpConnection) Close() {

}

// GetRemoteIp get remote ip.
func (t *UdpConnection) GetRemoteIp() string {
	return ""
}

// GetRemotePort get remote port.
func (t *UdpConnection) GetRemotePort() int {
	return 0
}

// GetRemoteAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *UdpConnection) GetRemoteAddress() string {
	return ""
}

// GetLocalIp get local ip.
func (t *UdpConnection) GetLocalIp() string {
	return ""
}

// GetLocalPort get local port.
func (t *UdpConnection) GetLocalPort() int {
	return 0
}

// GetLocalAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *UdpConnection) GetLocalAddress() string {
	return ""
}
