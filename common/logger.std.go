package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// std output to stdout, stdErr output to stderr.
var std = logrus.New()
var stdErr = logrus.New()

// formatter formats the output format.
type formatter struct {
	isStdout    bool
	serviceName string
}

// Format the input log.
func (f *formatter) Format(e *logrus.Entry) ([]byte, error) {
	// Implode the data to string with `k=v` format.
	dataString := ""
	if len(e.Data) != 0 {
		for k, v := range e.Data {
			dataString += fmt.Sprintf("%s=%+v ", k, v)
		}
		// Trim the trailing whitespace.
		dataString = dataString[0 : len(dataString)-1]
	}
	// Get service name.
	name := f.serviceName
	// Level like: DEBUG, INFO, WARN, ERROR, FATAL.
	level := strings.ToUpper(e.Level.String())
	// Get the time with YYYY-mm-dd H:i:s format.
	time := e.Time.Format("2006-01-02 15:04:05")
	// Get the message.
	msg := e.Message

	// Set the color of the level with STDOUT.
	stdLevel := ""
	switch level {
	case "DEBUG":
		stdLevel = color.New(color.FgWhite).Sprint(level)
	case "TRACE":
		stdLevel = color.New(color.FgWhite).Sprint(level)
	case "INFO":
		stdLevel = color.New(color.FgCyan).Sprint(" " + level)
	case "WARN":
		stdLevel = color.New(color.FgYellow).Sprint(" " + level)
	case "ERROR":
		stdLevel = color.New(color.FgRed).Sprint(level)
	case "FATAL":
		stdLevel = color.New(color.FgHiRed).Sprint(level)
	}

	body := fmt.Sprintf("[%s] %5s %s %s %s", name, stdLevel, time, REQUEST_ID, msg)
	data := fmt.Sprintf(" (%s)", dataString)

	// Hide the data if there's no data.
	if len(e.Data) == 0 {
		data = ""
	}

	// Mix the body and the data.
	output := fmt.Sprintf("%s%s\n", body, data)

	return []byte(output), nil
}

// LogInit Init initializes the global logger.
func LogInit(c *cli.Context) {
	var stdFmt logrus.Formatter

	// Create the formatter for stdout output.
	stdFmt = &formatter{
		isStdout:    false,
		serviceName: "Control",
	}

	// Std logger.
	std.Out = os.Stdout
	std.Formatter = stdFmt

	// StdErr logger
	stdErr.Formatter = stdFmt

	switch strings.ToUpper(c.String("log-level")) {
	case "FATAL":
		std.Level = logrus.FatalLevel
		stdErr.Level = logrus.FatalLevel
	case "ERROR":
		std.Level = logrus.ErrorLevel
		stdErr.Level = logrus.ErrorLevel
	case "WARN":
		std.Level = logrus.WarnLevel
		stdErr.Level = logrus.WarnLevel
	case "INFO":
		std.Level = logrus.InfoLevel
		stdErr.Level = logrus.InfoLevel
	case "DEBUG":
		std.Level = logrus.DebugLevel
		stdErr.Level = logrus.DebugLevel
	case "TRACE":
		std.Level = logrus.TraceLevel
		stdErr.Level = logrus.TraceLevel
	default:
		std.Level = logrus.DebugLevel
		stdErr.Level = logrus.DebugLevel
	}
}

func LogDebug(msg interface{}) {
	message("Debug", msg)
}
func LogTrace(msg interface{}) {
	message("Debug", msg)
}
func LogInfo(msg interface{}) {
	message("Info", msg)
}
func LogWarn(msg interface{}) {
	message("Warn", msg)
}
func LogError(msg interface{}) {
	message("Error", msg)
}
func LogFatal(msg interface{}) {
	message("Fatal", msg)
}

func LogDebugf(msg string, fds logrus.Fields) {
	fields("Debug", msg, fds)
}
func LogTracef(msg string, fds logrus.Fields) {
	fields("Trace", msg, fds)
}
func LogInfof(msg string, fds logrus.Fields) {
	fields("Info", msg, fds)
}
func LogWarnf(msg string, fds logrus.Fields) {
	fields("Warn", msg, fds)
}
func LogErrorf(msg string, fds logrus.Fields) {
	fields("Error", msg, fds)
}
func LogFatalf(msg string, fds logrus.Fields) {
	fields("Fatal", msg, fds)
}

func fields(lvl string, msg string, fds logrus.Fields) {
	s := std.WithFields(fds)
	sr := stdErr.WithFields(fds)

	switch lvl {
	case "Debug":
		s.Debug(msg)
	case "Info":
		s.Info(msg)
	case "Warn":
		s.Warn(msg)
	case "Error":
		sr.Error(msg)
	case "Fatal":
		sr.Fatal(msg)
	case "Trace":
		s.Trace(msg)
	}
}

func message(lvl string, msg interface{}) {
	switch lvl {
	case "Debug":
		std.Debug(msg)
	case "Info":
		std.Info(msg)
	case "Warn":
		std.Warn(msg)
	case "Error":
		stdErr.Error(msg)
	case "Fatal":
		stdErr.Fatal(msg)
	case "Trace":
		std.Trace(msg)
	}
}
