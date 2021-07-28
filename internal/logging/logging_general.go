package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	traceLog     *log.Logger
	debugLog     *log.Logger
	infoLog      *log.Logger
	warningLog   *log.Logger
	errorLog     *log.Logger
	logLock      sync.Mutex
	theLogHeader string
)

const (
	tagTrace   = "TRACE: "
	tagDebug   = "DEBUG: "
	tagInfo    = "INFO : "
	tagWarning = "WARN : "
	tagError   = "ERROR: "
)

const (
	flags = log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC
)

type logLevel int

// Log levels, from most verbose (LogLevelTrace) to least (LogLevelDisabled)
const (
	LogLevelTrace = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelDisabled
)

/*
* log.Output() is called from safelyLogContent() in this package.
* log.Output() uses depth to determine the stack frame from which
* it collects the source File and source Line number of the log call.
*
* Stack frames when calling that log.Output() function:
*
* Original Caller					--4
* 	=> Warn(f)/Info(f)/Trace(f)/Debug(f)/Error(f)	--3
*		=> safelyLog(f)/safelyLogWithName(f)	--2
*			=> safelyLogContent()		--1
*
* So, depth parameter value is set to 4
* For detailed info: https://golang.org/pkg/log/#Output
 */
const originalCallerRelativeDepth = 4

var pid string
var savedLogLevel logLevel
var logger *log.Logger

const logFlags = log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC

func init() {
	pid = fmt.Sprintf("[%d] ", os.Getpid())
	savedLogLevel = LogLevelInfo
	logger = log.New(os.Stderr, "", logFlags)
}

var logLevelMap = map[string]logLevel{
	"TRACE":    LogLevelTrace,
	"DEBUG":    LogLevelDebug,
	"INFO":     LogLevelInfo,
	"WARN":     LogLevelWarn,
	"ERROR":    LogLevelError,
	"DISABLED": LogLevelDisabled,
}

var logLevelStrings = map[logLevel]string{
	LogLevelTrace:    "TRACE",
	LogLevelDebug:    "DEBUG",
	LogLevelInfo:     "INFO",
	LogLevelWarn:     "WARN",
	LogLevelError:    "ERROR",
	LogLevelDisabled: "DISABLED",
}

func curLogLevel() logLevel {
	return savedLogLevel
}

// Errors related to logging
var (
	ErrUnknownLogLevel = errors.New("Unknown log level")
)

// SetLogLevel sets the current log level.  Input parameter must be one of
// TRACE, DEBUG, INFO, WARN, ERROR, DISABLED
func SetLogLevel(newLevel string) error {

	level, ok := logLevelMap[newLevel]
	if !ok {
		// illegel log level
		return ErrUnknownLogLevel
	}
	savedLogLevel = level
	return nil
}

// GetLogLevel returns the current log level
func GetLogLevel() string {
	return logLevelStrings[savedLogLevel]
}

// SetLogger specifies to use a different logger when logging LRPC related messages.
// By default, logging is sent to stderr
func SetLogger(lg *log.Logger) {
	logger = lg
}

// GetLogger returns the logger used in logging LRPC related messages
func GetLogger() *log.Logger {
	return logger
}

// Shutdown redirects logging of LRPC related messages back to stderr.  Usually not needed.
func Shutdown() {
	logLock.Lock()
	defer logLock.Unlock()
	logger.SetOutput(os.Stderr)
}

func safelyLogContent(lvl logLevel, depth int, content string) {
	if curLogLevel() > lvl { // Only log if current log level requires it
		return
	}
	tag := logLevelStrings[lvl]
	logLock.Lock()
	defer logLock.Unlock()
	logger.Output(depth, tag+" "+pid+content)
}

func safelyLog(lvl logLevel, depth int, v ...interface{}) {
	safelyLogContent(lvl, depth, fmt.Sprint(v...))
}

func safelyLogf(lvl logLevel, depth int, format string, v ...interface{}) {
	safelyLogContent(lvl, depth, fmt.Sprintf(format, v...))
}

// Info logs a message as Info level message
func Info(v ...interface{}) {
	safelyLog(LogLevelInfo, originalCallerRelativeDepth, v...)
}

func Infof(format string, v ...interface{}) {
	safelyLogf(LogLevelInfo, originalCallerRelativeDepth, format, v...)
}

// Debug logs a message as Debug level message
func Debug(v ...interface{}) {
	safelyLog(LogLevelDebug, originalCallerRelativeDepth, v...)
}

func Debugf(format string, v ...interface{}) {
	safelyLogf(LogLevelDebug, originalCallerRelativeDepth, format, v...)
}

// Warn logs a message as Warn level message
func Warn(v ...interface{}) {
	safelyLog(LogLevelWarn, originalCallerRelativeDepth, v...)
}

func Warnf(format string, v ...interface{}) {
	safelyLogf(LogLevelWarn, originalCallerRelativeDepth, format, v...)
}

// Error logs a message as Error level message
func Error(v ...interface{}) {
	safelyLog(LogLevelError, originalCallerRelativeDepth, v...)
}

func Errorf(format string, v ...interface{}) {
	safelyLogf(LogLevelError, originalCallerRelativeDepth, format, v...)
}

// Trace logs a message as Trace level message
func Trace(v ...interface{}) {
	safelyLog(LogLevelTrace, originalCallerRelativeDepth, v...)
}

func Tracef(format string, v ...interface{}) {
	safelyLogf(LogLevelTrace, originalCallerRelativeDepth, format, v...)
}
