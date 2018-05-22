package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/test/assert"
)

func trapStderr(t *testing.T, f func(stderr *bufio.Reader)) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not open temporary stderr replacement pipe: %s", err.Error())
		return
	}
	defer writer.Close()
	defer reader.Close()

	save := os.Stderr
	defer func() {
		os.Stderr = save
		resetLoggers()
	}()
	os.Stderr = writer

	resetLoggers()
	f(bufio.NewReader(reader))
}

func testConsumeConsoleLogs(t *testing.T, level string, get func() *log.Logger) {
	lineRegex := func(msg string) string {
		return fmt.Sprintf(
			`^\[%s\] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} %s$`,
			level,
			msg,
		)
	}

	trapStderr(t, func(stderr *bufio.Reader) {
		get().Println("test stderr")
		testLine, err := stderr.ReadString('\n')
		assert.Nil(t, err)
		assert.MatchesRegex(t, strings.TrimSpace(testLine), lineRegex("test stderr"))

		ch, restore := ConsumeConsoleLogs(10)
		get().Println("test channel")

		msg := <-ch

		assert.Equal(t, msg.Level, level)
		assert.MatchesRegex(t, strings.TrimSpace(msg.Message), lineRegex("test channel"))

		restore()

		get().Println("test stderr after")

		testLine, err = stderr.ReadString('\n')
		assert.Nil(t, err)
		assert.MatchesRegex(t, strings.TrimSpace(testLine), lineRegex("test stderr after"))

		assert.ChannelEmpty(t, ch)
	})
}

func TestConsoleConsoleLogsInfo(t *testing.T) {
	testConsumeConsoleLogs(t, "info", Info)
}

func TestConsoleConsoleLogsError(t *testing.T) {
	testConsumeConsoleLogs(t, "error", Error)
}

func TestConsoleConsoleLogsDebug(t *testing.T) {
	savedChoice := logLevelChoice.Choice
	defer func() {
		logLevelChoice.Choice = savedChoice
	}()
	logLevelChoice.Choice = ptr.String("debug")

	testConsumeConsoleLogs(t, "debug", Debug)
}
