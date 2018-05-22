package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	saveFlags := logFlags
	defer func() {
		os.Stderr = save
		logFlags = saveFlags
		resetLoggers()
	}()
	os.Stderr = writer
	logFlags = 0

	resetLoggers()
	f(bufio.NewReader(reader))
}

func testConsumeConsoleLogs(t *testing.T, level string, get func() *log.Logger) {
	expectedLine := func(msg string) string {
		return fmt.Sprintf(`[%s] %s`, level, msg) + "\n"
	}

	trapStderr(t, func(stderr *bufio.Reader) {
		get().Println("test stderr")
		testLine, err := stderr.ReadString('\n')
		assert.Nil(t, err)
		assert.Equal(t, testLine, expectedLine("test stderr"))

		ch, restore := ConsumeConsoleLogs(10)
		get().Println("test channel")

		msg := <-ch

		assert.Equal(t, msg.Level, level)
		assert.Equal(t, msg.Message, expectedLine("test channel"))

		restore()

		get().Println("test stderr after")

		testLine, err = stderr.ReadString('\n')
		assert.Nil(t, err)
		assert.Equal(t, testLine, expectedLine("test stderr after"))

		assert.ChannelEmpty(t, ch)
	})
}

func TestConsumeConsoleLogsInfo(t *testing.T) {
	testConsumeConsoleLogs(t, "info", Info)
}

func TestConsumeConsoleLogsError(t *testing.T) {
	testConsumeConsoleLogs(t, "error", Error)
}

func TestConsumeConsoleLogsDebug(t *testing.T) {
	savedChoice := logLevelChoice.Choice
	defer func() {
		logLevelChoice.Choice = savedChoice
	}()
	logLevelChoice.Choice = ptr.String("debug")

	testConsumeConsoleLogs(t, "debug", Debug)
}
