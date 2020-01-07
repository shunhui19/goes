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
	OnMessage func(connection ConnectionInterface, data []byte)
	// OnError emitted when a error occurs with connection.
	OnError func(code int, msg string)
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
	// byteWrite bytes written.
	byteWrite int
	// currentPackageLength current package length.
	currentPackageLength int
}

// Send send data on the connection.
// if connection is closing or closed, return.
// if connection have not been established, then save encode data info application send buffer,
// otherwise direct send encode of data into system socket send buffer.
//
// return value.
// if return true, indicate the message already send to operating system layer socket send buffer.
// else return false, send fail.
// return nil indicate the message already send to application send buffer, waiting for into operating system socket send buffer.
func (t *TcpConnection) Send(buffer string, raw bool) interface{} {
	if t.status == StatusClosing || t.status == StatusClosed {
		return false
	}

	// try to call protocol::encode(send_buffer) before sending into the application send buffer.
	if raw == false && t.Protocol != nil {
		buffer = string(t.Protocol.Encode([]byte(buffer)))
		if len(buffer) == 0 {
			return nil
		}
	}

	// when the connection have not been established, save encode data into send buffer.
	if t.status != StatusEstablished {
		if len(t.sendBuffer) > 0 && t.bufferIsFull() {
			t.baseConnection.SendFail++
			return false
		}
		// the encode of data into application send buffer.
		t.sendBuffer = append(t.sendBuffer, []byte(buffer)...)
		t.checkBufferWillFull()
		return nil
	}

	// when the connection is established, send data directed.
	if len(t.sendBuffer) == 0 {
		n, err := (*t.socket).Write([]byte(buffer))
		if err != nil {
			lib.Warn("send data error: %v", err.Error())
			return false
		}
		// send success.
		if n == len(buffer) {
			t.byteWrite += n
			return true
		}
		// send only part of the data.
		if n > 0 {
			t.sendBuffer = []byte(buffer[n:])
			t.byteWrite += n
		} else {
			// connection whether is closed.
			if t.socket == nil {
				t.baseConnection.SendFail++
				if t.OnError != nil {
					t.OnError(2, "client closed!")
				}
				t.destroy()
				return false
			}
			// the fail of data send again.
			t.sendBuffer = []byte(buffer)
			t.checkBufferWillFull()
			t.write()
			return nil
		}
	} else {
		if t.bufferIsFull() {
			t.baseConnection.SendFail++
			return false
		}
		t.sendBuffer = append(t.sendBuffer, []byte(buffer)...)
		t.checkBufferWillFull()
	}
	return nil
}

// bufferIsFull whether send buffer is full.
func (t *TcpConnection) bufferIsFull() bool {
	if len(t.sendBuffer) >= MaxSendBufferSize {
		if t.OnError != nil {
			lib.Warn("code: %d, msg: %s", 2, "send buffer full and drop package")
			t.OnError(2, "msg:send buffer full and drop package")
		}
		return true
	}
	return false
}

// checkBufferWillFull check whether the send buffer will be full.
func (t *TcpConnection) checkBufferWillFull() {
	if len(t.sendBuffer) >= MaxSendBufferSize {
		if t.OnBuffFull != nil {
			t.OnBuffFull()
		}
	}
}

// Close close connection.
func (t *TcpConnection) Close(data string) {
	if t.status == StatusClosing || t.status == StatusClosed {
		return
	} else {
		if data != "" {
			t.Send(data, false)
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

// GetSendBufferQueueSize get send buffer queue size.
func (t *TcpConnection) GetSendBufferQueueSize() int {
	return len(t.sendBuffer)
}

// GetRecvBufferQueueSize get recv buffer queue size.
func (t *TcpConnection) GetRecvBufferQueueSize() int {
	return len(t.recvBuffer)
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
			t.Close("server close client")
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

// write write data into socket.
func (t *TcpConnection) write() {
	n, _ := (*t.socket).Write(t.sendBuffer)
	if n == len(t.sendBuffer) {
		t.byteWrite += n
		t.sendBuffer = t.sendBuffer[:0]
		if t.OnBufferDrain != nil {
			t.OnBufferDrain()
		}
		if t.status == StatusClosing {
			t.destroy()
		}
		return
	}

	if n > 0 {
		t.byteWrite += n
		t.sendBuffer = t.sendBuffer[n:]
	} else {
		t.baseConnection.SendFail++
		t.destroy()
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
