// updConnection.
package connections

import (
	"goes/lib"
	"goes/protocols"
	"net"
	"strconv"
	"strings"
)

// UpdConnection struct.
type UdpConnection struct {
	baseConnection BaseConnection
	// Protocol Application layer protocol.
	Protocol protocols.Protocol
	// socket udp socket.
	socket net.PacketConn
	// remoteAddress remote address.
	remoteAddress net.Addr
}

// Send send data on the connection.
func (u *UdpConnection) Send(buffer string, raw bool) interface{} {
	if raw == false && u.Protocol != nil {
		buffer := u.Protocol.Encode([]byte(buffer))
		if len(buffer) == 0 {
			return nil
		}
	}

	n, err := (u.socket).WriteTo([]byte(buffer), u.remoteAddress)
	if err != nil {
		lib.Warn("udp write to client error:", err)
		return false
	}
	return len(buffer) == n
}

// Close close connection.
func (u *UdpConnection) Close(data string) {
	if len(data) != 0 {
		u.Send(data, false)
	}
	return
}

// GetRemoteIp get remote ip.
func (u *UdpConnection) GetRemoteIp() string {
	return strings.Split(u.remoteAddress.String(), ":")[0]
}

// GetRemotePort get remote port.
func (u *UdpConnection) GetRemotePort() int {
	port, _ := strconv.Atoi(strings.Split(u.remoteAddress.String(), ":")[1])
	return port
}

// GetRemoteAddress get remote address, the format is like this 127.0.0.1:8080.
func (u *UdpConnection) GetRemoteAddress() string {
	return u.remoteAddress.String()
}

// GetLocalIp get local ip.
func (u *UdpConnection) GetLocalIp() string {
	return strings.Split((u.socket).LocalAddr().String(), ":")[0]
}

// GetLocalPort get local port.
func (u *UdpConnection) GetLocalPort() int {
	port, _ := strconv.Atoi(strings.Split((u.socket).LocalAddr().String(), ":")[1])
	return port
}

// GetLocalAddress get remote address, the format is like this http://127.0.0.1:8080.
func (u *UdpConnection) GetLocalAddress() string {
	return (u.socket).LocalAddr().String()
}

// AddRequestCount increase request record.
func (u *UdpConnection) AddRequestCount() {
	u.baseConnection.TotalRequest++
}

func NewUdpConnection(socket net.PacketConn, remoteAddr net.Addr) *UdpConnection {
	udp := &UdpConnection{}
	udp.socket = socket
	udp.remoteAddress = remoteAddr
	return udp
}
