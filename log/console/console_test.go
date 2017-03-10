package console

import (
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type testCase struct {
	level             string
	expectDebugLogger bool
	expectInfoLogger  bool
	expectErrorLogger bool
}

func (tc *testCase) run(t *testing.T) {
	savedChoice := logLevelChoice.Choice
	defer func() {
		logLevelChoice.Choice = savedChoice
	}()

	logLevelChoice.Choice = &tc.level

	assert.Group(fmt.Sprintf("level = %q", tc.level), t, func(g *assert.G) {
		if tc.expectDebugLogger {
			assert.NotSameInstance(t, Debug(), nullLogger)
			assert.SameInstance(t, Debug(), debugLogger)
		} else {
			assert.SameInstance(t, Debug(), nullLogger)
			assert.NotSameInstance(t, Debug(), debugLogger)
		}

		if tc.expectInfoLogger {
			assert.NotSameInstance(t, Info(), nullLogger)
			assert.SameInstance(t, Info(), infoLogger)
		} else {
			assert.SameInstance(t, Info(), nullLogger)
			assert.NotSameInstance(t, Info(), infoLogger)
		}

		if tc.expectErrorLogger {
			assert.NotSameInstance(t, Error(), nullLogger)
			assert.SameInstance(t, Error(), errorLogger)
		} else {
			assert.SameInstance(t, Error(), nullLogger)
			assert.NotSameInstance(t, Error(), errorLogger)
		}
	})
}

func TestConsoleLoggers(t *testing.T) {
	testCases := []testCase{
		{
			level: "none",
		},
		{
			level:             "error",
			expectErrorLogger: true,
		},
		{
			level:             "info",
			expectErrorLogger: true,
			expectInfoLogger:  true,
		},
		{
			level:             "debug",
			expectErrorLogger: true,
			expectInfoLogger:  true,
			expectDebugLogger: true,
		},
		{
			level:             "unexpected-value",
			expectErrorLogger: true,
			expectInfoLogger:  true,
		},
	}

	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestUninitializedChoice(t *testing.T) {
	savedChoice := logLevelChoice.Choice
	defer func() {
		logLevelChoice.Choice = savedChoice
	}()

	logLevelChoice.Choice = nil

	assert.Equal(t, logLevel(), defaultOrdinal)
}
