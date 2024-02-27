package logger

import (
	"fmt"
	"time"
)

// type check
var _ ILogger = (*ConsoleLogger)(nil)

func NewConsoleLogger(bufferSize uint) ILogger {
	return &ConsoleLogger{
		logsChannel: make(chan *Message, bufferSize),
	}
}

func (cl *ConsoleLogger) Start() {
	go func() {
		for msg := range cl.logsChannel {
			fmt.Printf("[%s] %s: %s\n", levelToString(msg.Level), msg.Time.Format(time.DateTime), msg.Message)
		}
	}()
}

func (cl *ConsoleLogger) Debug(msg string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.logsChannel <- makeMessage(Debug, msg)
}

func (cl *ConsoleLogger) Info(msg string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.logsChannel <- makeMessage(Info, msg)
}

func (cl *ConsoleLogger) Warn(msg string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.logsChannel <- makeMessage(Warn, msg)
}

func (cl *ConsoleLogger) Error(msg string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.logsChannel <- makeMessage(Error, msg)
}
