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

type Level int

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
var useColor = true

func SetLevel(l Level) {
	level = l
}

func DisableColor() {
	useColor = false
}

func colorize(l Level, s string) string {
	if !useColor {
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

func Debug(v ...interface{}) error {
	return output(os.Stdout, LevelDebug, v...)
}

func Debugf(format string, v ...interface{}) error {
	return outputf(os.Stdout, LevelDebug, format, v...)
}

func Info(v ...interface{}) error {
	return output(os.Stdout, LevelInfo, v...)
}

func Infof(format string, v ...interface{}) error {
	return outputf(os.Stdout, LevelInfo, format, v...)
}

func Warn(v ...interface{}) error {
	return output(os.Stdout, LevelWarn, v...)
}

func Warnf(format string, v ...interface{}) error {
	return outputf(os.Stdout, LevelWarn, format, v...)
}

func Error(v ...interface{}) error {
	return output(os.Stderr, LevelError, v...)
}

func Errorf(format string, v ...interface{}) error {
	return outputf(os.Stderr, LevelError, format, v...)
}
