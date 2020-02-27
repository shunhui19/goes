// updConnection.
package connections

import (
	"net"
	"strconv"

	"github.com/shunhui19/goes/lib"
	"github.com/shunhui19/goes/protocols"
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
}

// Send send data on the connection.
// if return true indicate message send success, otherwise return false indicate send fail.
// return nil indicate no data to send.
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

	err := u.socket.Close()
	if err != nil {
		lib.Warn(err.Error())
	}
}

// GetRemoteIP get remote ip.
func (u *UDPConnection) GetRemoteIP() string {
	IP, _, _ := net.SplitHostPort(u.remoteAddress.String())
	return IP
}

// GetRemotePort get remote port.
func (u *UDPConnection) GetRemotePort() int {
	_, port, _ := net.SplitHostPort(u.remoteAddress.String())
	p, _ := strconv.Atoi(port)
	return p
}

// GetRemoteAddress get remote address, the format is like this 127.0.0.1:8080.
func (u *UDPConnection) GetRemoteAddress() string {
	return u.remoteAddress.String()
}

// GetLocalIP get local ip.
func (u *UDPConnection) GetLocalIP() string {
	IP, _, _ := net.SplitHostPort((u.socket).LocalAddr().String())
	return IP
}

// GetLocalPort get local port.
func (u *UDPConnection) GetLocalPort() int {
	_, port, _ := net.SplitHostPort((u.socket).LocalAddr().String())
	p, _ := strconv.Atoi(port)
	return p
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
	return udp
}
