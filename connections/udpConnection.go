// updConnection.
package connections

import (
	"goes/lib"
	"goes/protocols"
	"net"
	"strconv"
	"strings"
)

// UDPConnection struct.
type UDPConnection struct {
	baseConnection BaseConnection
	// Protocol Application layer protocol.
	Protocol protocols.Protocol
	// socket udp socket.
	socket net.PacketConn
	// remoteAddress remote address.
	remoteAddress net.Addr
	// MaxPackageSize set the maximum packet size for receive and send.
	MaxPackageSize int
}

// Send send data on the connection.
func (u *UDPConnection) Send(buffer string, raw bool) interface{} {
	if raw == false && u.Protocol != nil {
		switch result := u.Protocol.Encode([]byte(buffer)).(type) {
		case []byte:
			buffer = string(result)
		case string:
			buffer = result
		}
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
func (u *UDPConnection) Close(data string) {
	if len(data) != 0 {
		u.Send(data, false)
	}
	return
}

// GetRemoteIP get remote ip.
func (u *UDPConnection) GetRemoteIP() string {
	return strings.Split(u.remoteAddress.String(), ":")[0]
}

// GetRemotePort get remote port.
func (u *UDPConnection) GetRemotePort() int {
	port, _ := strconv.Atoi(strings.Split(u.remoteAddress.String(), ":")[1])
	return port
}

// GetRemoteAddress get remote address, the format is like this 127.0.0.1:8080.
func (u *UDPConnection) GetRemoteAddress() string {
	return u.remoteAddress.String()
}

// GetLocalIP get local ip.
func (u *UDPConnection) GetLocalIP() string {
	return strings.Split((u.socket).LocalAddr().String(), ":")[0]
}

// GetLocalPort get local port.
func (u *UDPConnection) GetLocalPort() int {
	port, _ := strconv.Atoi(strings.Split((u.socket).LocalAddr().String(), ":")[1])
	return port
}

// GetLocalAddress get remote address, the format is like this http://127.0.0.1:8080.
func (u *UDPConnection) GetLocalAddress() string {
	return (u.socket).LocalAddr().String()
}

// AddRequestCount increase request record.
func (u *UDPConnection) AddRequestCount() {
	u.baseConnection.TotalRequest++
}

// NewUDPConnection new a object of UDPConnection.
func NewUDPConnection(socket net.PacketConn, remoteAddr net.Addr) *UDPConnection {
	udp := &UDPConnection{}
	udp.socket = socket
	udp.remoteAddress = remoteAddr
	udp.MaxPackageSize = DefaultMaxPackageSize
	return udp
}
