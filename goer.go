package goes

import (
	"goes/lib"
	"net"
)

const (
	// Version the version of goes.
	Version = 0.1

	// Status the status of starting.
	StatusStarting = 1
	// Status the status of running.
	StatusRunning = 2
	// Status the status of shutdown.
	StatusShutdown = 4
	// Status the status of reloading.
	StatusReloading = 8
)

// build-in protocols.
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
	// Transport the protocol of layer.
	Transport string
	// Protocol the protocol of application.
	Protocol string
	// Daemon daemon start.
	Daemon bool
	// LogFile log file.
	LogFile string
	// GracefulStop graceful stop or not.
	gracefulStop bool
	// mainSocket listening socket.
	mainSocket net.Listener
	// socketName socket name, the format is like this http://127.0.0.1:8080
	socketName string
	// OnConnect emitted when a socket connection is successfully established.
	OnConnect func()
	// OnMessage emitted when data is received.
	OnMessage func()
	// OnClose emitted when other end of the socket sends a FIN packet.
	OnClose func()
	// OnError emitted when an error occurs with connection.
	OnError func()
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
	g.initWorkers()
	g.installSignal()
	g.saveMainPid()
	g.displayUI()
	g.resetStd()
	g.monitorWorkers()
}

// checkEnv check environment.
func (g *Goer) checkEnv() {

}

// init.
func (g *Goer) init() {

}

// parseCommand parse command.
func (g *Goer) parseCommand() {

}

// daemon run as daemon mode.
func (g *Goer) daemon() {

}

// initWorkers init all worker instances.
func (g *Goer) initWorkers() {

}

// saveMainPid save pid.
func (g *Goer) saveMainPid() {

}

// displayUI display starting UI.
func (g *Goer) displayUI() {

}

// resetStd redirect standard input and output.
func (g *Goer) resetStd() {

}

// monitorWorkers monitor all child goroutine.
func (g *Goer) monitorWorkers() {

}

// installSignal install signal handler.
func (g *Goer) installSignal() {

}

// NewGoer create object of Goer with socketName.
func NewGoer(socketName string) *Goer {
	if socketName == "" {
		lib.Fatal("the socket address can not be empty")
	}
	return &Goer{socketName: socketName}
}
