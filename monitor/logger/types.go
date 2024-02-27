package logger

import (
	"sync"
	"time"
)

type LogLevel uint8

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

type Message struct {
	Level   LogLevel
	Message string
	Time    time.Time
}

type ILogger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Start()
}

type ConsoleLogger struct {
	logsChannel chan *Message
	mutex       sync.Mutex
}
