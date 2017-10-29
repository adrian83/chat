package logger

import (
	"fmt"
	"os"
	"time"

	logrus "github.com/sirupsen/logrus"
)

const (
	lInfo    = "INFO"
	lDebug   = "DEBUG"
	lWarning = "WARNING"
	lError   = "ERROR"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	logrus.SetOutput(os.Stderr)
}

// Log is a struct representing single log entry.
type Log struct {
	StructName   string
	FunctionName string
	Timestamp    time.Time
	Message      string
	Level        string
}

// New returns new Log struct filled with given data.
func New(strName, funcName, msg, level string) *Log {
	return &Log{
		StructName:   strName,
		FunctionName: funcName,
		Message:      msg,
		Level:        level,
		Timestamp:    time.Now().UTC(),
	}
}

// String returns string representation of the Log struct.
func (l *Log) String() string {
	return fmt.Sprintf("[%s] [%7s] [%20s] [%20s] %s", l.Timestamp.Format(timeFormat), l.Level, l.StructName, l.FunctionName, l.Message)
}

// Info logs message with Info level.
func Info(strName, funcName, msg string) {
	logrus.Info(New(strName, funcName, msg, lInfo).String())
}

// Infof logs parametrized message with Info level.
func Infof(strName, funcName, msgTmpl string, params ...interface{}) {
	Info(strName, funcName, fmt.Sprintf(msgTmpl, params...))
}

// Debug logs message with Debug level.
func Debug(strName, funcName, msg string) {
	logrus.Debug(New(strName, funcName, msg, lDebug).String())
}

// Debugf logs parametrized message with Debug level.
func Debugf(strName, funcName, msgTmpl string, params ...interface{}) {
	Debug(strName, funcName, fmt.Sprintf(msgTmpl, params...))
}

// Warn logs message with Warn level.
func Warn(strName, funcName, msg string) {
	logrus.Warn(New(strName, funcName, msg, lWarning).String())
}

// Warnf logs parametrized message with Warn level.
func Warnf(strName, funcName, msgTmpl string, params ...interface{}) {
	Warn(strName, funcName, fmt.Sprintf(msgTmpl, params...))
}

// Error logs message with Error level.
func Error(strName, funcName, msg string) {
	logrus.Error(New(strName, funcName, msg, lError).String())
}

// Errorf logs parametrized message with Error level.
func Errorf(strName, funcName, msgTmpl string, params ...interface{}) {
	Error(strName, funcName, fmt.Sprintf(msgTmpl, params...))
}
