// TCPConnection.
// note: in OSX system, the default max concurrency is limited to 128,
// so when your concurrency number more than 128 will occur error: connect: connection reset by peer.
// you should by setting: 'sysctl -w kern.ipc.somaxconn=xxx', xxx is the max currency number.
// link https://github.com/golang/go/issues/20960.
package connections

import (
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/shunhui19/goes/lib"
	"github.com/shunhui19/goes/protocols"
)

const (
	// ReadBufferSize read buffer size.
	ReadBufferSize = 65535
	// MaxSendBufferSize maximum size of send buffer.
	MaxSendBufferSize = 104856
	// DefaultMaxSendBufferSize default maximum size of send buffer.
	DefaultMaxSendBufferSize = 104856
	// MaxPackageSize maximum size of acceptable buffer.
	MaxPackageSize = 104856
	// DefaultMaxPackageSize default maximum size of acceptable buffer.
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

	// GoerSendFail send fail.
	GoerSendFail = 2
)

// TCPConnection struct.
type TCPConnection struct {
	// baseConnection.
	baseConnection BaseConnection
	// OnMessage emitted when data is received.
	OnMessage func(connection Connection, data []byte)
	// OnError emitted when a error occurs with connection.
	OnError func(connection Connection, code int, message string)
	// OnClose emitted when the other end of the socket send a FIN package.
	OnClose func(connection Connection)
	// OnBufferFull emitted when the send buffer becomes full.
	OnBufferFull func(connection Connection)
	// OnBufferDrain emitted when the send buffer becomes empty.
	OnBufferDrain func(connection Connection)
	// Transport transport.
	Transport string
	// Protocol application layer protocol.
	Protocol protocols.Protocol
	// ID the id of connection.
	ID int
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
	// Connections all connection instances, key is connection id and value is *net.Conn.
	Connections CStore
	// status connection status.
	status int
	// byteRead bytes read.
	byteRead int
	// byteWrite bytes written.
	byteWrite int
	// currentPackageLength current package length.
	currentPackageLength int
}

// Send send data on the connection.
// if connection is closing or closed, return false.
// if connection have not been established, then save encode data info application send buffer,
// otherwise direct send encode of data into system socket send buffer.
//
// return value.
// if return true, indicate the message already send to operating system layer socket send buffer.
// else return false, send fail.
// return nil indicate the message already send to application send buffer, waiting for into operating system socket send buffer.
func (t *TCPConnection) Send(buffer string, raw bool) interface{} {
	if t.status == StatusClosing || t.status == StatusClosed {
		return false
	}

	// try to call protocol::encode(send_buffer) before sending into the application send buffer.
	// every protocol maybe return different type after encode func.
	if raw == false && t.Protocol != nil {
		switch result := t.Protocol.Encode([]byte(buffer)).(type) {
		case []byte:
			buffer = string(result)
		case string:
			buffer = result
		}
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
		// connection maybe closed.
		if err != nil {
			if t.socket == nil || err == io.EOF {
				t.baseConnection.SendFail++
				if t.OnError != nil {
					t.OnError(t, GoerSendFail, "client is closed")
				}
				t.destroy()
				return false
			}
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
			if t.socket == nil || err == io.EOF {
				t.baseConnection.SendFail++
				if t.OnError != nil {
					t.OnError(t, GoerSendFail, "client is closed")
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
func (t *TCPConnection) bufferIsFull() bool {
	if len(t.sendBuffer) >= t.MaxSendBufferSize {
		if t.OnError != nil {
			t.OnError(t, GoerSendFail, "send buffer full and drop package")
		}
		return true
	}
	return false
}

// checkBufferWillFull check whether the send buffer will be full.
func (t *TCPConnection) checkBufferWillFull() {
	if len(t.sendBuffer) >= t.MaxSendBufferSize {
		if t.OnBufferFull != nil {
			t.OnBufferFull(t)
		}
	}
}

// Close close connection.
func (t *TCPConnection) Close(data string) {
	if t.status == StatusClosing || t.status == StatusClosed {
		return
	}

	if data != "" {
		t.Send(data, false)
	}
	t.status = StatusClosing
	t.destroy()
}

// GetRemoteIP get remote ip.
func (t *TCPConnection) GetRemoteIP() string {
	return strings.Split(t.remoteAddress, ":")[0]
}

// GetRemotePort get remote port.
func (t *TCPConnection) GetRemotePort() int {
	if t.remoteAddress != "" {
		port, _ := strconv.Atoi(strings.Split(t.remoteAddress, ":")[1])
		return port
	}
	return 0
}

// GetRemoteAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *TCPConnection) GetRemoteAddress() string {
	return t.remoteAddress
}

// GetLocalIP get local ip.
func (t *TCPConnection) GetLocalIP() string {
	return strings.Split((*t.socket).LocalAddr().String(), ":")[0]
}

// GetLocalPort get local port.
func (t *TCPConnection) GetLocalPort() int {
	addr := strings.Split((*t.socket).LocalAddr().String(), ":")
	port, _ := strconv.Atoi(addr[1])
	return port
}

// GetLocalAddress get remote address, the format is like this http://127.0.0.1:8080.
func (t *TCPConnection) GetLocalAddress() string {
	return (*t.socket).LocalAddr().String()
}

// GetSendBufferQueueSize get send buffer queue size.
func (t *TCPConnection) GetSendBufferQueueSize() int {
	return len(t.sendBuffer)
}

// GetRecvBufferQueueSize get recv buffer queue size.
func (t *TCPConnection) GetRecvBufferQueueSize() int {
	return len(t.recvBuffer)
}

// Read read data from socket.
func (t *TCPConnection) Read() {
	// ssl handle.
	if t.Transport == "ssl" {

	}

	bytesPool := lib.NewBytesPool(1024)
	for {
	READ:
		buf := bytesPool.Get()
		n, err := (*t.socket).Read(buf.B)
		if err != nil || err == io.EOF {
			t.Close("server close client")
			return
		}

		// check connection closed.
		if n == 0 {
			t.destroy()
			return
		}

		t.byteRead += n
		t.recvBuffer = append(t.recvBuffer, buf.B[:n]...)
		bytesPool.Put(buf)
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
					input := t.Protocol.Input(t.recvBuffer, t.MaxPackageSize)
					switch input.(type) {
					case int:
						t.currentPackageLength = input.(int)
						// the package length is unknown.
						if t.currentPackageLength == 0 {
							goto READ
						} else if t.currentPackageLength > 0 && t.currentPackageLength < t.MaxPackageSize {
							// data is not enough for a package.
							if t.currentPackageLength > t.byteRead {
								goto READ
							}
						}
					case bool:
						if input == false {
							lib.Warn("error package. the package of current connection maxsize is: %d", t.MaxPackageSize)
							t.destroy()
							return
						}
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
func (t *TCPConnection) write() {
	n, _ := (*t.socket).Write(t.sendBuffer)
	if n == len(t.sendBuffer) {
		t.byteWrite += n
		t.sendBuffer = t.sendBuffer[:0]
		if t.OnBufferDrain != nil {
			t.OnBufferDrain(t)
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
func (t *TCPConnection) destroy() {
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
		t.OnClose(t)
	}

	if t.status == StatusClosed {
		t.Connections.Del(t.ID)
		t.OnMessage, t.OnError, t.OnClose, t.OnBufferDrain, t.OnBufferFull = nil, nil, nil, nil, nil
	}
}

// NewTCPConnection new a tcp connection.
func NewTCPConnection(socket *net.Conn, remoteAddress string) *TCPConnection {
	tcp := &TCPConnection{}
	tcp.baseConnection.ConnectionCount++
	tcp.socket = socket
	tcp.MaxSendBufferSize = DefaultMaxSendBufferSize
	tcp.MaxPackageSize = DefaultMaxPackageSize
	tcp.recvBuffer = make([]byte, 0, ReadBufferSize)
	tcp.status = StatusEstablished
	tcp.remoteAddress = remoteAddress
	return tcp
}
