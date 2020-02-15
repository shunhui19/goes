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
	"time"
)

const (
	// Version the version of goes.
	Version = 0.1
	// MaxUDPPackageSize max udp package size.
	MaxUDPPackageSize = 65536

	// StatusStarting the status of starting.
	StatusStarting = 1
	// StatusRunning the status of running.
	StatusRunning = 2
	// StatusShutdown the status of shutdown.
	StatusShutdown = 4
	// StatusReloading the status of reloading.
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

// Goer the main-goroutine server.
type Goer struct {
	// Name the name of main goroutine.
	//Name string
	// User unix user of process, needs appropriate privileges, usually root.
	//User string
	// Reloadable reloadable.
	//Reloadable bool
	// ReusePort reuse port.
	//ReusePort bool
	// Transport the protocol of transport layer, if transport layer protocol is empty,
	// the default protocol is tcp.
	Transport string
	// Protocol the protocol of application layer, the type is interface of protocol,
	// if no set, the default protocol is tcp.
	Protocol protocols.Protocol
	// Daemon daemon start.
	Daemon bool
	// isForked whether is a fork process after by os/exec.Command().
	isForked bool
	// StdoutFile stdout file.
	StdoutFile string
	// LogFile log file.
	LogFile string
	// PidFile pid file.
	PidFile string
	// mainSocket listening socket.
	mainSocket interface{}
	// socketName socket name, the format is like this http://127.0.0.1:8080
	socketName string
	// mainPid the pid of the socket process.
	mainPid int
	// rootPath root path.
	rootPath string
	// Connections store all Connections of client.
	Connections sync.Map
	// gracefulWait wait for connections graceful exit which belongs to old process.
	gracefulWait *sync.WaitGroup
	// connectionID unique connection id.
	connectionID int
	// status current status.
	status int
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
	//OnMainGoroutineReload func()
	// OnMainGoroutineStop emitted when the main goroutine terminated.
	//OnMainGoroutineStop func()
}

// RunAll start server.
func (g *Goer) RunAll() {
	g.checkEnv()
	g.init()
	g.parseCommand()
	g.daemon()
	g.resetStd()
	g.listen()
	g.installSignal()
	//g.displayUI()
	//g.monitorWorkers()
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

	g.status = StatusStarting
}

// parseCommand parse command.
func (g *Goer) parseCommand() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: yourExecuteFile <command> [mode]\ncommand:\nstart\tStart goer in DEBUG mode.\n\tUse mode -d to start in DAEMON mode.\nstop\tStop goer.\nreload\tReload codes.\n")
		os.Exit(0)
	}

	// parse command.
	command := strings.Trim(os.Args[1], " ")
	command2 := ""
	if os.Args[1] == "start" {
		model := "debug"
		if len(os.Args) == 3 {
			command2 = os.Args[2]
			model = "daemon"
		}
		lib.Info("Goer start in %s mode.", strings.ToUpper(model))
	}

	switch command {
	// daemon main goroutine.
	// use os/exec.Command() function to launch a child process of parent program.
	case "main":
		g.Daemon = true
		g.isForked = true
	case "start":
		if _, err := os.Stat(g.PidFile); err == nil {
			lib.Fatal("Already running or pid: %s file exist", g.PidFile)
		}

		if command2 == "-d" {
			g.Daemon = true
		}

		g.mainPid = os.Getpid()
		g.saveMainPid()
	case "stop":
		processPid, err := g.getPid()
		if err != nil {
			lib.Fatal(err.Error())
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
	// reload.
	case "reload":
		processPid, err := g.getPid()
		if err != nil {
			lib.Fatal(err.Error())
		}
		process, err := os.FindProcess(processPid)
		if err != nil {
			lib.Fatal("Unable to find process ID[%v]", processPid)
		}

		// send SIGHUP signal to process.
		err = process.Signal(syscall.SIGQUIT)
		if err != nil {
			lib.Fatal("Send SIGQUIT fail, error: %v", err)
		}
		//time.Sleep(1 * time.Second)

		os.Exit(0)
	default:
		lib.Fatal("Unknown command: %v", os.Args[1])
	}
}

// daemon run as daemon mode.
func (g *Goer) daemon() {
	if !g.Daemon {
		return
	}

	if !g.isForked {
		cmd := exec.Command(os.Args[0], "main")
		cmd.Start()
		g.mainPid = cmd.Process.Pid
		g.saveMainPid()
		os.Exit(0)
	}
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
// the defined signal is:
// SIGINT => stop
// SIGTERM => graceful stop
// SIGUSR1 => reload
// SIGQUIT => graceful reload
// SIGUSR2 => status
func (g *Goer) installSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for signalType := range ch {
		switch signalType {
		// stop process in debug mode with Ctrl+c.
		case syscall.SIGINT:
			g.stopAll(ch, signalType)
		// kill signal in bash shell.
		case syscall.SIGKILL | syscall.SIGTERM:
			g.stopAll(ch, signalType)
		// graceful reload
		case syscall.SIGQUIT:
			signal.Stop(ch)
			g.reload()
			os.Exit(0)
		}
	}
}

// reload graceful to restart service.
// copy parent socket file descriptor to fork a child process.
func (g *Goer) reload() {
	g.status = StatusReloading

	// notice parents process stop accept new connection.
	err := g.mainSocket.(*net.TCPListener).SetDeadline(time.Now())
	if err != nil {
		lib.Fatal("listener socket set timeout fail: %v", err)
	}
	lib.Info("Goer is reloading...")

	// emitted when goer process get reload signal.
	if g.OnGoerReload != nil {
		g.OnGoerReload()
	}

	// get parent process of listener file descriptor.
	f, err := g.mainSocket.(*net.TCPListener).File()
	if err != nil {
		lib.Fatal("ListenFD error: %v", err)
	}

	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), f.Fd()},
	}
	// avoid exec reload command too long when exec many times.
	if len(os.Args) < 3 {
		os.Args = append(os.Args, "graceful")
	}
	childProcessPid, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		lib.Fatal("Fail to fork: %v", err)
	}

	// write the child process FD to PidFile.
	os.Remove(g.PidFile)
	g.mainPid = childProcessPid
	g.saveMainPid()
	lib.Info("Received SIGQUIT to fork-exec: %v", childProcessPid)

	// wait for the connections finished which belongs to parent process,
	// and then parent process exit.
	// if the connections finished time Exceeded maximum completion time,
	// the server will close all client.
	if err := g.gracefulWaitTimeout(time.Minute); err != nil {
		g.Connections.Range(func(key, value interface{}) bool {
			value.(connections.ConnectionInterface).Close("server is graceful restart")
			return true
		})
		lib.Fatal("Timeout when graceful")
	}
	lib.Info("Stop parent process success")
}

// stopAll stop.
func (g *Goer) stopAll(ch chan os.Signal, sig os.Signal) {
	g.status = StatusShutdown
	lib.Info("Goer is stopping...")
	signal.Stop(ch)
	lib.Info("Received signal type: %v", sig)

	// execute exit.
	g.stop()

	// remove pid file.
	err := os.Remove(g.PidFile)
	if err != nil {
		lib.Fatal("Remove pid file fail: %v", err)
	}
	lib.Info("Goer stop success")

	os.Exit(0)
}

// stop stop goer instance.
func (g *Goer) stop() {
	// emitted OnGoerStop callback func.
	if g.OnGoerStop != nil {
		g.OnGoerStop()
	}

	// close client.
	g.Connections.Range(func(k, connection interface{}) bool {
		connection.(connections.ConnectionInterface).Close("")
		return true
	})
}

// listen create a listen socket.
func (g *Goer) listen() {
	if g.socketName == "" {
		return
	}

	if g.mainSocket == nil {
		switch g.Transport {
		case "tcp", "tcp4", "tcp6", "unix", "unixpacket", "ssl":
			if len(os.Args) > 2 && os.Args[2] == "graceful" {
				file := os.NewFile(3, "")
				listener, err := net.FileListener(file)
				if err != nil {
					lib.Fatal("Fail to listen tcp: %v", err)
				}
				g.mainSocket = listener.(*net.TCPListener)
			} else {
				addr, err := net.ResolveTCPAddr(g.Transport, g.socketName)
				if err != nil {
					lib.Fatal("fail to resolve addr: %v", err)
				}
				listener, err := net.ListenTCP("tcp", addr)
				if err != nil {
					lib.Fatal("fail to listen tcp: %v", err)
				}
				g.mainSocket = listener
			}
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
		g.status = StatusRunning

		go g.resumeAccept()
	}
}

// resumeAccept accept new Connections.
func (g *Goer) resumeAccept() {
	switch g.Transport {
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket", "ssl":
		g.acceptTCPConnection()
	case "udp", "upd4", "udp6", "unixgram":
		g.acceptUDPConnection()
	}
}

// acceptTCPConnection accept a tcp connection.
func (g *Goer) acceptTCPConnection() {
	for {
		newSocket, err := g.mainSocket.(*net.TCPListener).Accept()
		if err != nil {
			// stop accept new connection when received reload signal.
			if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
				lib.Info("parent process stop accept new connection")
				return
			}
			lib.Error("unAccept client socket, reason: %s", err.Error())
			continue
		}
		connection := connections.NewTCPConnection(&newSocket, newSocket.RemoteAddr().String())
		connection.Transport = g.Transport
		connection.Protocol = g.Protocol
		connection.OnMessage = g.OnMessage
		connection.OnClose = g.OnClose
		connection.OnError = g.OnError
		connection.OnBufferDrain = g.OnBufferDrain
		connection.OnBuffFull = g.OnBufferFull
		// store all client connection.
		g.Connections.Store(g.generateConnectionID(), connection)

		// trigger OnConnect if is set.
		if g.OnConnect != nil {
			g.OnConnect(connection)
		}

		g.gracefulWait.Add(1)
		go func() {
			defer connection.Close("")
			connection.Read()
			// waiting for reload signal and one by one close old client.
			g.gracefulWait.Done()
		}()
	}
}

// acceptUDPConnection accept a udp package.
func (g *Goer) acceptUDPConnection() {
	for {
		recvBuffer := make([]byte, MaxUDPPackageSize)
		n, addr, err := g.mainSocket.(net.PacketConn).ReadFrom(recvBuffer)
		if err != nil {
			lib.Warn("ReadFrom the %d data of udp error: %v", n, err)
			return
		}
		go func() {
			connection := connections.NewUDPConnection(g.mainSocket.(net.PacketConn), addr)
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

// RemoveConnection delete connection from Connections store.
func (g *Goer) RemoveConnection(connectionID int) {
	g.Connections.Delete(connectionID)
}

// NewGoer create object of Goer with socketName, application layer protocol and transport layer protocol,
// if applicationProtocol is empty.
func NewGoer(socketName string, applicationProtocol protocols.Protocol, transportProtocol string) *Goer {
	if socketName == "" {
		lib.Fatal("the socket address can not be empty")
	}

	return &Goer{socketName: socketName, Protocol: applicationProtocol, Transport: transportProtocol, gracefulWait: &sync.WaitGroup{}}
}

// generateConnectionID generate unique connection id.
func (g *Goer) generateConnectionID() int {
	maxUnsignedInt := int(2147483647)
	if g.connectionID >= maxUnsignedInt {
		g.connectionID = 1
	}
	for g.connectionID < maxUnsignedInt {
		// start from 1.
		if g.connectionID == 0 {
			g.connectionID++
			continue
		}
		// judge current id whether has used.
		if _, ok := g.Connections.Load(g.connectionID); !ok {
			break
		}
		g.connectionID++
	}
	return g.connectionID
}

// getPid get the pid value from PidFle.
func (g *Goer) getPid() (int, error) {
	if _, err := os.Stat(g.PidFile); err == nil {
		data, err := ioutil.ReadFile(g.PidFile)
		if err != nil {
			lib.Fatal("Goer not run.")
		}
		processPid, err := strconv.Atoi(string(data))
		if err != nil {
			lib.Fatal("Unable to read and parse process pid found in: %v", g.PidFile)
		}
		return processPid, nil
	}

	return 0, fmt.Errorf("goer not run")
}

// gracefulWaitTimeout set graceful timeout.
func (g *Goer) gracefulWaitTimeout(duration time.Duration) error {
	timeout := time.NewTicker(duration)
	wait := make(chan struct{})

	go func() {
		g.gracefulWait.Wait()
		wait <- struct{}{}
	}()

	select {
	case <-timeout.C:
		return fmt.Errorf("time out")
	case <-wait:
		return nil
	}
}
