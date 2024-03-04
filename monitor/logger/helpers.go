package logger

import (
	"time"
)

// helper for boilerplate reduction
func makeMessage(level LogLevel, msg string) *Message {
	return &Message{
		Level:   level,
		Message: msg,
		Time:    time.Now(),
	}
}

// converts a log level to its string representation
func levelToString(level LogLevel) string {
	switch level {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	default:
		return "error"
	}
}
