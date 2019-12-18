// the library of logging, can use global method(as Info, Debug, Warn, Error, Panic, Fatal) everywhere,
// also define a struct of Log implement the interface of Logger
package lib

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// InfoLevel the default level of logging.
	InfoLevel = iota + 1
	// DebugLevel the debug level of logging, usually use in development.
	DebugLevel
	// WarnLevel log the type of warn message.
	WarnLevel
	// ErrorLevel log the type of error message.
	ErrorLevel
	// PanicLevel log a message and panic.
	PanicLevel
	// FatalLevel log a message and call os.Exit(1).
	FatalLevel
)

// levelTagMap set tag of log.
var levelTagMap = map[int]string{
	// unused color
	InfoLevel: "INFO",
	// use color of white
	DebugLevel: White("DEBUG"),
	// use color of yellow
	WarnLevel: Yellow("WARN"),
	// use color of magenta
	ErrorLevel: Magenta("ERROR"),
	// use color of magenta
	PanicLevel: Magenta("PANIC"),
	// use color of red
	FatalLevel: Red("FATAL"),
}

type Logger interface {
	// Info log the type of info message.
	Info(format string, a ...interface{})
	// Debug log the type of Debug message.
	Debug(format string, a ...interface{})
	// Warn log the type of warn message.
	Warn(format string, a ...interface{})
	// Error log the type of error message.
	Error(format string, a ...interface{})
	// Panic log a message and call panic.
	Panic(format string, a ...interface{})
	// Fatal log a message and call os.Exit(1).
	Fatal(format string, a ...interface{})
}

// Info log the type of info message.
func Info(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), InfoLevel)
}

// Debug log the type of Debug message.
func Debug(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), DebugLevel)
}

// Warn log the type of warn message.
func Warn(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), WarnLevel)
}

// Error log the type of error message.
func Error(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), ErrorLevel)
}

// Panic log a message and call panic.
func Panic(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), ErrorLevel)
	panic(fmt.Sprintf(format, a...))
}

// Fatal log a message and call os.Exit(1).
func Fatal(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), FatalLevel)
	os.Exit(1)
}

// outPrint echo content of message to Stdout.
func outPrint(msg string, level int) {
	logHead := "[" + time.Now().Format("2006-01-02 15:04:05.000") + "]"
	logHead += "[" + levelTagMap[level] + "] "
	logHead += callName(2) + " "
	fmt.Println(logHead + msg)
}

// Log struct of Log
type Log struct {
}

func (l *Log) Info(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), InfoLevel)
}

func (l *Log) Debug(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), DebugLevel)
}

func (l *Log) Warn(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), WarnLevel)
}

func (l *Log) Error(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), ErrorLevel)
}

func (l *Log) Panic(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), PanicLevel)
	panic(fmt.Sprintf(format, a...))
}

func (l *Log) Fatal(format string, a ...interface{}) {
	outPrint(fmt.Sprintf(format, a...), FatalLevel)
	os.Exit(1)
}

// callStack get call information with file, line, function.
func callName(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	return file[strings.LastIndex(file, "/src/")+5:] + ":" + strconv.Itoa(line) + " [" + name + "]"
}
