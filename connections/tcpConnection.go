// TcpConnection.
package connections

import (
	"goes/lib"
	"goes/protocols"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	// ReadBufferSize read buffer size.
	ReadBufferSize = 65535
	// MaxSendBufferSize maximum size of send buffer.
	MaxSendBufferSize = 104856
	// MaxSendBufferSize default maximum size of send buffer.
	DefaultMaxSendBufferSize = 104856
	// MaxPackageSize maximum size of acceptable buffer.
	MaxPackageSize = 104856
	// MaxPackageSize default maximum size of acceptable buffer.
	DefaultMaxPackageSize = 1048560

	// StatusInitial initial status of connection.
	StatusInitial = 0
	// StatusConnecting connecting status of connection.
	StatusConnecting = 1
	// StatusEstablished established status of connection.
	StatusEstablished = 2
	// StatusClosing closing status of connection.
	StatusClosing = 4
	// StatusClosed closed status of connection.
	StatusClosed = 8
)

// TcpConnection struct.
type TcpConnection struct {
	// baseConnection.
	baseConnection BaseConnection
	// OnMessage emitted when data is received.
	OnMessage func(connection *TcpConnection, data []byte)
	// OnError emitted when a error occurs with connection.
	OnError func()
	// OnClose emitted when the other end of the socket send a FIN package.
	OnClose func()
	// OnBuffFull emitted when the send buffer becomes full.
	OnBuffFull func()
	// OnBufferDrain emitted when the send buffer becomes empty.
	OnBufferDrain func()
	// Transport transport.
	Transport string
	// Protocol application layer protocol.
	Protocol protocols.Protocol
	// Id the id of connection.
	Id int
	// idRecorder id recorder.
	idRecorder int
	// MaxSendBufferSize set the maximum send buffer size for the current connection.
	MaxSendBufferSize int
	// MaxPackageSize set the maximum acceptable packet size for the current connection.
	MaxPackageSize int
	// socket tcp socket.
	socket *net.Conn
	// remoteAddress remote address.
	remoteAddress string
	// recvBuffer receive buffer.
	recvBuffer []byte
	// sendBuffer send buffer.
	sendBuffer []byte
	// connections all connection instances, key is connection id and value is *net.Conn.
	connections sync.Map
	// status connection status.
	status int
	// Goer which goer belong to.
	//Goer *goes.Goer
	// byteRead bytes read.
	byteRead int
	// currentPackageLength current package length.
	currentPackageLength int
}

// Send send data on the connection.
func (t *TcpConnection) Send(data string) {

}

// Close close connection.
func (t *TcpConnection) Close(data string) {
	if t.status == StatusClosing || t.status == StatusClosed {
		return
	} else {
		if data != "" {
			t.Send(data)
		}
		t.status = StatusClosing
	}
	t.destroy()
	//if len(t.sendBuffer) == 0 {
	//} else {
	//
	//}
}

// GetRemoteIp get remote ip.
func (t *TcpConnection) GetRemoteIp() string {
	return strings.Split(t.remoteAddress, ":")[0]
}

// GetRemotePort get remote port.
func (t *TcpConnection) GetRemotePort() int {
	if t.remoteAddress != "" {
		port, _ := strconv.Atoi(strings.Split(t.remoteAddress, ":")[1])
		return port
	}
	return 0
}

// GetRemoteAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *TcpConnection) GetRemoteAddress() string {
	return t.remoteAddress
}

// GetLocalIp get local ip.
func (t *TcpConnection) GetLocalIp() string {
	return strings.Split((*t.socket).LocalAddr().String(), ":")[0]
}

// GetLocalPort get local port.
func (t *TcpConnection) GetLocalPort() int {
	addr := strings.Split((*t.socket).LocalAddr().String(), ":")
	port, _ := strconv.Atoi(addr[1])
	return port
}

// GetLocalAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *TcpConnection) GetLocalAddress() string {
	return (*t.socket).LocalAddr().String()
}

// Read read data from socket.
func (t *TcpConnection) Read() {
	// ssl handle.
	if t.Transport == "ssl" {

	}

	for {
	READ:
		// determine the size of receive buf in every package.
		//err := binary.Read(*t.socket, binary.BigEndian, &size)
		//if err != nil {
		//	lib.Warn("determine the size error: %v", err.Error())
		//}
		buf := make([]byte, 1024)

		n, err := (*t.socket).Read(buf)
		if err != nil && err != io.EOF {
			lib.Warn(err.Error())
			return
		}

		// check connection closed.
		if n == 0 {
			t.destroy()
			return
		} else {
			t.byteRead += n
			t.recvBuffer = append(t.recvBuffer, buf[:n]...)
		}

		// if the application layer protocol has been set up.
		if t.Protocol != nil {
			for len(t.recvBuffer) > 0 {
				// the current packet length is known.
				if t.currentPackageLength > 0 {
					// data is not enough for a package.
					if t.currentPackageLength > t.byteRead {
						goto READ
					}
				} else {
					// get length of package from protocol interface.
					input := t.Protocol.Input(t.recvBuffer)
					switch input.(type) {
					case int:
						t.currentPackageLength = input.(int)
					case bool:
						if input == false {
							t.Close("")
							return
						}
					default:
					}

					// the package length is unknown.
					if t.currentPackageLength == 0 {
						goto READ
					} else if t.currentPackageLength > 0 && t.currentPackageLength < MaxPackageSize {
						// data is not enough for a package.
						if t.currentPackageLength > t.byteRead {
							goto READ
						}
					} else {
						lib.Warn("error package. package_length=%d", t.currentPackageLength)
						t.destroy()
						return
					}
				}

				// the data is enough for a package.
				t.baseConnection.TotalRequest++
				// get a full package from the buffer.
				oneRequestBuffer := t.recvBuffer[:t.currentPackageLength]
				// remove the current package from the receive buffer.
				t.recvBuffer = t.recvBuffer[t.currentPackageLength:]
				// reset the current package length.
				t.currentPackageLength = 0
				if t.OnMessage == nil {
					continue
				}

				// decode request buffer before emitted OnMessage func.
				t.OnMessage(t, t.Protocol.Decode(oneRequestBuffer))
			}
		}

		if len(t.recvBuffer) == 0 {
			continue
		}

		// application protocol is not set.
		t.baseConnection.TotalRequest++
		if t.OnMessage == nil {
			t.recvBuffer = t.recvBuffer[:0]
			t.byteRead = 0
			continue
		}
		t.OnMessage(t, t.recvBuffer)
		t.recvBuffer = t.recvBuffer[:0]
	}
}

// destroy destroy connection.
func (t *TcpConnection) destroy() {
	if t.status == StatusClosed {
		return
	}

	// close socket.
	err := (*t.socket).Close()
	if err != nil {
		lib.Warn(err.Error())
	}
	t.status = StatusClosed

	// trigger OnClose func.
	if t.OnClose != nil {
		t.OnClose()
	}

	// trigger OnClose func of protocol.

	// whether gc ???.
	if t.status == StatusClosed {
		// remove from goer.connections.
		//if t.Goer != nil {
		//	t.Goer.RemoveConnection(t.Id)
		//}
	}
}

// NewTcpConnection new a tcp connection.
func NewTcpConnection(socket *net.Conn, remoteAddress string) *TcpConnection {
	tcp := &TcpConnection{}
	tcp.baseConnection.ConnectionCount++
	tcp.Id++
	tcp.idRecorder++
	if tcp.idRecorder == math.MaxInt32 {
		tcp.idRecorder = 0
	}
	tcp.socket = socket
	tcp.MaxSendBufferSize = DefaultMaxSendBufferSize
	tcp.MaxPackageSize = DefaultMaxPackageSize
	tcp.recvBuffer = make([]byte, 0, ReadBufferSize)
	tcp.status = StatusEstablished
	tcp.remoteAddress = remoteAddress
	//tcp.connections.Store(tcp.Id, tcp)
	return tcp
}
