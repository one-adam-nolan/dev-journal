package logs

import (
	"github.com/gookit/color"
)

type DjLogger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type GooKitLogger struct {
}

func (l *GooKitLogger) Debug(msg string) {
	color.Debug.Prompt(msg)
	color.Debug.Tips(msg)
}

func (l *GooKitLogger) Info(msg string) {
	color.Info.Prompt(msg)
	color.Info.Tips(msg)

}

func (l *GooKitLogger) Warn(msg string) {
	color.Warn.Prompt(msg)
	color.Warn.Tips(msg)
}

func (l *GooKitLogger) Error(msg string) {
	color.Error.Prompt(msg)
	color.Error.Tips(msg)
}


func GetColorfulLogger() DjLogger {
	return &GooKitLogger{}
}