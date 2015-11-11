// Package log provides functions for logging debug, informational, warning,
// and error messages to standard output/error. Clients should set the current
// log level; only messages at that level or higher will actually be logged.
// Compared to Go's standard log package, this package supports colored output.
//
// Inspired by https://github.com/cloudflare/cfssl/blob/master/log/log.go
package log

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/colorstring"
)

// The Level type is the type of all log levels.
type Level int

// The different log levels.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelPrefix = [...]string{
	LevelDebug: "DEBUG: ",
	LevelInfo:  "INFO: ",
	LevelWarn:  "WARNING: ",
	LevelError: "ERROR: ",
}

var levelColor = [...]string{
	LevelDebug: "[blue]",
	LevelInfo:  "[cyan]",
	LevelWarn:  "[yellow]",
	LevelError: "[red]",
}

var level = LevelDebug

// UseColor enables colorized output is set to true.
var UseColor = true

// SetLevel changes the current log level to l.
func SetLevel(l Level) {
	level = l
}

func colorize(l Level, s string) string {
	if !UseColor {
		return s
	}
	return colorstring.Color(levelColor[l] + s)
}

func output(w io.Writer, l Level, v ...interface{}) error {
	if l < level {
		return nil
	}
	_, err := fmt.Fprint(w, colorize(l, levelPrefix[l]+fmt.Sprintln(v...)))
	return err
}

func outputf(w io.Writer, l Level, format string, v ...interface{}) error {
	if l < level {
		return nil
	}
	_, err := fmt.Fprintf(w, colorize(l, levelPrefix[l]+format), v...)
	return err
}

// Debug logs a debug message to stdout.
func Debug(v ...interface{}) error {
	return output(os.Stdout, LevelDebug, v...)
}

// Debugf logs a formatted debug message to stdout.
func Debugf(format string, v ...interface{}) error {
	return outputf(os.Stdout, LevelDebug, format, v...)
}

// Info logs an informational message to stdout.
func Info(v ...interface{}) error {
	return output(os.Stdout, LevelInfo, v...)
}

// Infof logs a formatted informational message to stdout.
func Infof(format string, v ...interface{}) error {
	return outputf(os.Stdout, LevelInfo, format, v...)
}

// Warn logs a warning message to stdout.
func Warn(v ...interface{}) error {
	return output(os.Stdout, LevelWarn, v...)
}

// Warnf logs a formatted warning message to stdout.
func Warnf(format string, v ...interface{}) error {
	return outputf(os.Stdout, LevelWarn, format, v...)
}

// Error logs an error message to stderr.
func Error(v ...interface{}) error {
	return output(os.Stderr, LevelError, v...)
}

// Errorf logs a formatted error message to stderr.
func Errorf(format string, v ...interface{}) error {
	return outputf(os.Stderr, LevelError, format, v...)
}
