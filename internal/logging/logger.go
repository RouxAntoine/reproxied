package logging

import (
	"fmt"
	"io"
	"os"
)

// Logger capability.
type Logger interface {
	Error(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

// Writer local alias for io.StringWriter.
type Writer io.StringWriter

type reProxiedLogger struct {
	prefix string
	level  Level
	writer Writer
}

// NewReProxiedLogger default logger constructor with info level and os.Stdout output.
func NewReProxiedLogger(pluginName string) Logger {
	return &reProxiedLogger{
		prefix: pluginName,
		level:  Levels.INFO,
		writer: os.Stdout,
	}
}

// NewReProxiedLoggerWithLevel logger constructor with custom level.
func NewReProxiedLoggerWithLevel(pluginName string, writer Writer, level Level) Logger {
	return &reProxiedLogger{
		prefix: pluginName,
		level:  level,
		writer: writer,
	}
}

// Error filter log with level upper or equals to error level.
func (rl *reProxiedLogger) Error(message string, args ...interface{}) {
	if rl.level <= Levels.ERR {
		rl.logColored(Levels.ERR, message, args)
	}
}

// Warn filter log with level upper or equals to warning level.
func (rl *reProxiedLogger) Warn(message string, args ...interface{}) {
	if rl.level <= Levels.WARN {
		rl.logColored(Levels.WARN, message, args)
	}
}

// Info filter log with level upper or equals to info level.
func (rl *reProxiedLogger) Info(message string, args ...interface{}) {
	if rl.level <= Levels.INFO {
		rl.logColored(Levels.INFO, message, args)
	}
}

// Debug filter log with level upper or equals to debug level.
func (rl *reProxiedLogger) Debug(message string, args ...interface{}) {
	if rl.level <= Levels.DEBUG {
		rl.logColored(Levels.DEBUG, message, args)
	}
}

func (rl *reProxiedLogger) logColored(level Level, message string, args []interface{}) {
	var coloredMessage string
	switch level {
	case Levels.ERR:
		coloredMessage = fmt.Sprintf("%s[ERR]%s %s", Color.RED, Color.CLEAR, message)
	case Levels.WARN:
		coloredMessage = fmt.Sprintf("%s[WARN]%s %s", Color.ORANGE, Color.CLEAR, message)
	case Levels.DEBUG:
		coloredMessage = fmt.Sprintf("%s[DEBUG]%s %s", Color.GREEN, Color.CLEAR, message)
	default:
		coloredMessage = fmt.Sprintf("%s[INFO]%s %s", Color.BLUE, Color.CLEAR, message)
	}

	rl.log(coloredMessage, args)
}

// log format a string message and write it to stdout.
func (rl *reProxiedLogger) log(message string, args []interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	decoratedMessage := fmt.Sprintf("%s[reproxied]%s - %s\n", Color.CYAN, Color.CLEAR, formattedMessage)
	_, _ = rl.writer.WriteString(decoratedMessage)
}
