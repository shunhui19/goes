package goes

import (
	"fmt"
	"goes/connections"
	"goes/lib"
	"goes/protocols"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	// Version the version of goes.
	Version = 0.1
	// MaxUdpPackageSize max udp package size.
	MaxUdpPackageSize = 65536

	// Status the status of starting.
	StatusStarting = 1
	// Status the status of running.
	StatusRunning = 2
	// Status the status of shutdown.
	StatusShutdown = 4
	// Status the status of reloading.
	StatusReloading = 8
)

// buildInTransports Go build-in transports protocols.
var buildInTransports = []map[string]string{
	{"tcp": "tcp"},
	{"tcp4": "tcp4"},
	{"tcp6": "tcp6"},
	{"udp": "udp"},
	{"unix": "unix"},
	{"unixpacket": "unixpacket"},
	{"ssl": "tcp"},
}

type Goer struct {
	// Name the name of main goroutine.
	Name string
	// User unix user of process, needs appropriate privileges, usually root.
	User string
	// Reloadable reloadable.
	Reloadable bool
	// ReusePort reuse port.
	ReusePort bool
	// Transport the protocol of transport layer, if transport layer protocol is empty,
	// the default protocol is tcp.
	Transport string
	// Protocol the protocol of application layer, the type is interface of protocol,
	// if no set, the default protocol is tcp.
	Protocol protocols.Protocol
	// Daemon daemon start.
	Daemon bool
	// StdoutFile stdout file.
	StdoutFile string
	// LogFile log file.
	LogFile string
	// PidFile pid file.
	PidFile string
	// GracefulStop graceful stop or not.
	gracefulStop bool
	// mainSocket listening socket.
	mainSocket interface{}
	// socketName socket name, the format is like this http://127.0.0.1:8080
	socketName string
	// mainPid the pid of the socket process.
	mainPid int
	// rootPath root path.
	rootPath string
	// connections store all connections of client.
	connections sync.Map
	// OnConnect emitted when a socket connection is successfully established.
	OnConnect func(connection connections.ConnectionInterface)
	// OnMessage emitted when data is received.
	OnMessage func(connection connections.ConnectionInterface, data []byte)
	// OnClose emitted when other end of the socket sends a FIN packet.
	OnClose func()
	// OnError emitted when an error occurs with connection.
	OnError func(code int, msg string)
	// OnBufferFull emitted when the send buffer becomes full.
	OnBufferFull func()
	// OnBufferDrain emitted when the send buffer is empty.
	OnBufferDrain func()
	// OnGoerStop emitted when goer process stop.
	OnGoerStop func()
	// OnGoerReload emitted when goer process get reload signal.
	OnGoerReload func()
	// OnMainGoroutineReload emitted when the main goroutine process get reload signal.
	OnMainGoroutineReload func()
	// OnMainGoroutineStop emitted when the main goroutine terminated.
	OnMainGoroutineStop func()
}

// RunAll start server.
func (g *Goer) RunAll() {
	g.checkEnv()
	g.init()
	g.parseCommand()
	g.daemon()
	g.resetStd()
	g.initWorkers()
	g.installSignal()
	g.saveMainPid()
	g.displayUI()
	g.monitorWorkers()
}

// checkEnv check environment.
func (g *Goer) checkEnv() {

}

// init.
func (g *Goer) init() {
	// check transport layer protocol.
	if g.Transport == "" {
		g.Transport = "tcp"
	}

	// get root path.
	g.rootPath, _ = os.Getwd()

	// default stdoutFile
	if g.StdoutFile == "" {
		g.StdoutFile = "/dev/null"
	}

	// default pid file.
	if g.PidFile == "" {
		_, prefix, _, _ := runtime.Caller(1)
		path := filepath.Dir(prefix)
		prefix = strings.ReplaceAll(prefix, "/", "_")
		subStr := strings.Split(prefix, ".")
		g.PidFile = path + "/" + subStr[0] + ".pid"
	}
}

// parseCommand parse command.
func (g *Goer) parseCommand() {
	if len(os.Args) < 2 {
		lib.Fatal("Usage: %s [start|stop] \n", os.Args[0])
	}

	// parse command.
	command := strings.Trim(os.Args[0], " ")
	command2 := ""
	if len(os.Args) == 3 {
		command2 = os.Args[2]
	}

	switch os.Args[1] {
	// main goroutine.
	case "main":
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGUSR2)
		go func() {
			// block waiting to receive signal.
			signalType := <-ch
			switch signalType {
			case syscall.SIGKILL | syscall.SIGTERM:
				fmt.Println(1111)
				// if receive one signal then stop receive others signal.
				signal.Stop(ch)
				lib.Info("Received signal type: %v", signalType)

				// remove pid file.
				err := os.Remove(g.PidFile)
				if err != nil {
					lib.Fatal("Remove pid file fail: %v", err)
				}

				os.Exit(0)
			case syscall.SIGUSR2:
				lib.Info("Receive user signal: %v", signalType)
			default:
			}
		}()
	case "start":
		if _, err := os.Stat(g.PidFile); err == nil {
			lib.Fatal("Already running or pid: %s file exist", g.PidFile)
		}

		if command2 == "-d" || g.Daemon {
			if g.Daemon == false {
				g.Daemon = true
			}

			cmd := exec.Command(command, "main")
			cmd.Start()
			lib.Info("Goer start in DAEMON mode.")
			g.mainPid = cmd.Process.Pid
			g.saveMainPid()
			lib.Info("Goer main socket process id: %v", g.mainPid)
			os.Exit(0)
		}
		g.mainPid = os.Getpid()
		g.saveMainPid()
		lib.Info("Goer start in DEBUG mode.")
	case "stop":
		if _, err := os.Stat(g.PidFile); err == nil {
			data, err := ioutil.ReadFile(g.PidFile)
			if err != nil {
				lib.Fatal("Goer not run.")
			}

			processPid, err := strconv.Atoi(string(data))
			if err != nil {
				lib.Fatal("Unable to read and parse process pid found in: %v", g.PidFile)
			}

			process, err := os.FindProcess(processPid)
			if err != nil {
				lib.Fatal("Unable to find process ID[%v]", processPid)
			}

			// remove pid file.
			os.Remove(g.PidFile)

			// kill process.
			lib.Info("Goer is stopping...")
			err = process.Kill()
			if err != nil {
				lib.Fatal("Goer stop fail, error: %v", err)
			}
			lib.Info("Goer stop success")

			os.Exit(0)
		} else {
			lib.Fatal("Goer not run.")
		}
	default:
		lib.Fatal("Unknown command: %v", os.Args[1])
	}
}

// daemon run as daemon mode.
func (g *Goer) daemon() {

}

// initWorkers init all worker instances.
func (g *Goer) initWorkers() {
	g.listen()
}

// saveMainPid save pid.
func (g *Goer) saveMainPid() {
	file, err := os.Create(g.PidFile)
	if err != nil {
		lib.Fatal("Unable to create pid file: %v\n", err)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(g.mainPid))
	if err != nil {
		lib.Fatal("Unable to write pid number to file: %v\n", err)
	}
	file.Sync()
}

// displayUI display starting UI.
func (g *Goer) displayUI() {

}

// resetStd redirect standard input and output.
func (g *Goer) resetStd() {
	if !g.Daemon {
		return
	}

	handle, err := os.OpenFile(g.StdoutFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		lib.Fatal("can not open StdoutFile: %v", g.StdoutFile)
	}
	os.Stdout = handle
	os.Stderr = handle
}

// monitorWorkers monitor all child goroutine.
func (g *Goer) monitorWorkers() {

}

// installSignal install signal handler.
func (g *Goer) installSignal() {

}

// listen create a listen socket.
func (g *Goer) listen() {
	if g.socketName == "" {
		return
	}

	if g.mainSocket == nil {
		switch g.Transport {
		case "tcp", "tcp4", "tcp6", "unix", "unixpacket", "ssl":
			listener, err := net.Listen(g.Transport, g.socketName)
			if err != nil {
				lib.Fatal(err.Error())
			}
			g.mainSocket = listener
		case "udp", "upd4", "udp6", "unixgram":
			listener, err := net.ListenPacket(g.Transport, g.socketName)
			if err != nil {
				lib.Fatal(err.Error())
			}
			g.mainSocket = listener
		default:
			lib.Fatal("unknown transport layer protocol")
		}

		lib.Info("server start success...")

		g.resumeAccept()
	}
}

// resumeAccept accept new connections.
func (g *Goer) resumeAccept() {
	switch g.Transport {
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket", "ssl":
		g.acceptTcpConnection()
	case "udp", "upd4", "udp6", "unixgram":
		g.acceptUdpConnection()
	}
}

// acceptTcpConnection accept a tcp connection.
func (g *Goer) acceptTcpConnection() {
	for {
		newSocket, err := g.mainSocket.(net.Listener).Accept()
		if err != nil {
			lib.Error("unAccept client:%s socket, reason: %s", newSocket.RemoteAddr().String(), err.Error())
			continue
		}
		connection := connections.NewTcpConnection(&newSocket, newSocket.RemoteAddr().String())
		// store all client connection.
		g.connections.Store(connection.Id, *connection)
		//connection.Goer = g
		connection.Transport = g.Transport
		connection.Protocol = g.Protocol
		connection.OnMessage = g.OnMessage
		connection.OnClose = g.OnClose
		connection.OnError = g.OnError
		connection.OnBufferDrain = g.OnBufferDrain
		connection.OnBuffFull = g.OnBufferFull

		// trigger OnConnect if is set.
		if g.OnConnect != nil {
			g.OnConnect(connection)
		}

		go func() {
			defer connection.Close("")
			connection.Read()
		}()
	}
}

// acceptUdpConnection accept a udp package.
func (g *Goer) acceptUdpConnection() {
	for {
		recvBuffer := make([]byte, MaxUdpPackageSize)
		n, addr, err := g.mainSocket.(net.PacketConn).ReadFrom(recvBuffer)
		if err != nil {
			lib.Warn("ReadFrom the %d data of udp error: %v", n, err)
			return
		}
		go func() {
			connection := connections.NewUdpConnection(g.mainSocket.(net.PacketConn), addr)
			connection.Protocol = g.Protocol
			if g.OnMessage != nil {
				if g.Protocol != nil {
					if n > 0 {
						input := g.Protocol.Input(recvBuffer)
						switch input.(type) {
						case int:
							if input.(int) == 0 {
								return
							}
							packet := recvBuffer[:input.(int)]
							recvBuffer = recvBuffer[input.(int):]
							data := g.Protocol.Decode(packet)
							g.OnMessage(connection, data)
						case bool:
							if input.(bool) == false {
								return
							}
						}
					}
				} else {
					g.OnMessage(connection, recvBuffer[:n])
				}
				connection.AddRequestCount()
			}
		}()
	}
}

// RemoveConnection delete connection from connections store.
func (g *Goer) RemoveConnection(connectionId int) {
	g.connections.Delete(connectionId)
}

// NewGoer create object of Goer with socketName, application layer protocol and transport layer protocol,
// if applicationProtocol is empty.
func NewGoer(socketName string, applicationProtocol protocols.Protocol, transportProtocol string) *Goer {
	if socketName == "" {
		lib.Fatal("the socket address can not be empty")
	}

	return &Goer{socketName: socketName, Protocol: applicationProtocol, Transport: transportProtocol}
}
