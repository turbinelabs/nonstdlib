package console

import (
	"io"
	"log"
	"os"
)

// LogMessage represents a console log level and message.
type LogMessage struct {
	Level   string
	Message string
}

// ConsumeConsoleLogs allows log messages to be temporarily diverted to a channel for
// testing. Messages created via the Info, Error, or Debug loggers are sent to the
// returned channel, which is created with the given capacity. If the channel is
// full, subsequent calls to any of the loggers will block. The returned function
// must be called to restore the default output. Nested calls of ConsumeConsoleLogs
// are not supported.
//
// This method is intended to aid in testing code that uses the console package and
// should not be used outside tests.
func ConsumeConsoleLogs(capacity int) (<-chan LogMessage, func()) {
	ch := make(chan LogMessage, capacity)

	set := func(log *log.Logger, w io.Writer, flags int) {
		log.SetOutput(w)
		log.SetFlags(flags)
	}

	set(errorLogger, &channelWriter{level: errorLevel, ch: ch}, 0)
	set(infoLogger, &channelWriter{level: infoLevel, ch: ch}, 0)
	set(debugLogger, &channelWriter{level: debugLevel, ch: ch}, 0)

	return ch, func() {
		set(errorLogger, os.Stderr, logFlags)
		set(infoLogger, os.Stderr, logFlags)
		set(debugLogger, os.Stderr, logFlags)
	}
}

type channelWriter struct {
	level string
	ch    chan LogMessage
}

func (w *channelWriter) Write(b []byte) (int, error) {
	w.ch <- LogMessage{
		Level:   w.level,
		Message: string(b),
	}
	return len(b), nil
}
