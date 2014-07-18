package log

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/colorstring"
)

const (
	LevelDebug = iota
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

var Level = LevelDebug
var UseColor = true

func colorize(l int, s string) string {
	if !UseColor {
		return s
	}
	return colorstring.Color(levelColor[l] + s)
}

func format(l int, v ...interface{}) string {
	msg := levelPrefix[l] + fmt.Sprintln(v...)
	return colorize(l, msg)
}

func output(w io.Writer, l int, v ...interface{}) error {
	if l >= Level {
		_, err := fmt.Fprint(w, format(l, v...))
		return err
	}
	return nil
}

func Debug(v ...interface{}) error {
	return output(os.Stdout, LevelDebug, v...)
}

func Info(v ...interface{}) error {
	return output(os.Stdout, LevelInfo, v...)
}

func Warn(v ...interface{}) error {
	return output(os.Stdout, LevelWarn, v...)
}

func Error(v ...interface{}) error {
	return output(os.Stderr, LevelError, v...)
}
