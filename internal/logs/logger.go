package logs

import (
	"github.com/gookit/color"
)

// Logger provides leveled log output.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

// ColorfulLogger implements Logger with terminal colors.
type ColorfulLogger struct{}

// NewColorfulLogger returns a colorful terminal logger.
func NewColorfulLogger() Logger {
	return &ColorfulLogger{}
}

func (l *ColorfulLogger) Debug(msg string) {
	color.Debug.Println(msg)
}

func (l *ColorfulLogger) Info(msg string) {
	color.Info.Println(msg)
}

func (l *ColorfulLogger) Warn(msg string) {
	color.Warn.Println(msg)
}

func (l *ColorfulLogger) Error(msg string) {
	color.Error.Println(msg)
}
