package logging

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/centrify/platform-go-sdk/testutils"
	"github.com/stretchr/testify/suite"
)

type LoggingTestSuite struct {
	testutils.CfyTestSuite
}

var goodLogLevels = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "DISABLED"}

func (s *LoggingTestSuite) TestGetSetLogLevel() {

	t := s.T()
	t.Log("Testing get/set log level")

	savedLogLevel := GetLogLevel()

	err := SetLogLevel("unknown")
	s.Assert().ErrorIs(err, ErrUnknownLogLevel, "Expect return error for unknown level")

	// test get/set of good levels
	for _, level := range goodLogLevels {
		err = SetLogLevel(level)
		s.Assert().NoError(err, "Should not encountered error in setting log level")
		s.Assert().Equal(level, GetLogLevel(), "GetLogLevel should return same value as level set")
	}

	SetLogLevel(savedLogLevel)
}

func (s *LoggingTestSuite) TestLogLevels() {
	t := s.T()
	a := s.Assert()

	t.Log("Test log messages at different level ")

	savedLogLevel := GetLogLevel()
	visibleMsg := "%s level messages should be visible when current log level is %s"
	hiddenMsg := "%s level messages should be hidden when current log level is %s"

	var err error
	var curlevel logLevel
	var msg string
	var msgLevel logLevel
	var hidden bool
	var logbuf bytes.Buffer
	var msgTag string

	// set up to use a new io.Writer for log
	savedLogger := GetLogger()
	newLogger := log.New(&logbuf, "", logFlags)
	SetLogger(newLogger)

	// generate random tag
	msgTag = fmt.Sprintf("Tag%v", rand.Int())

	for _, lvl := range goodLogLevels {
		// write messages for each level
		err = SetLogLevel(lvl)
		s.Assert().NoError(err, "Should not encountered error in setting log level")
		t.Logf("** Setting log level to %s", lvl)
		curlevel = logLevelMap[lvl]

		for _, msglvl := range goodLogLevels {
			logbuf.Reset() // reset test output buffer
			msgLevel = logLevelMap[msglvl]
			if msgLevel >= curlevel {
				msg = visibleMsg
				hidden = false
			} else {
				msg = hiddenMsg
				hidden = true
			}
			switch msgLevel {
			case LogLevelTrace:
				Trace("Trace level message: " + msgTag)
			case LogLevelDebug:
				Debug("Debug level message: " + msgTag)
			case LogLevelInfo:
				Info("Info level message: " + msgTag)
			case LogLevelWarn:
				Warn("Warn level message: " + msgTag)
			case LogLevelError:
				Error("Error level message: " + msgTag)
			case LogLevelDisabled:
				hidden = true
			default:
				t.Error("Should never reach here")
			}
			t.Logf(msg, msglvl, lvl)

			// verify data
			logmsg := logbuf.String()
			if hidden {
				a.Emptyf(logmsg, "There should be nothing logged when log message should be hidden for message level %s and current level %s",
					msglvl, lvl)
			} else {
				a.NotEmptyf(logmsg, "There should be nothing logged when log message should be hidden for message level %s and current level %s",
					msglvl, lvl)
				// check for log message tag
				a.Contains(logmsg, msgTag, "log message must contain tag")
				a.Contains(logmsg, msglvl, "log message must contain level")
			}

		}

	}
	// restore log configurations
	SetLogLevel(savedLogLevel)
	SetLogger(savedLogger)

}
func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}
