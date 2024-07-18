package logger

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"github.com/dubbogo/gost/log/logger"
	"gopkg.inshopline.com/commons/logx"
	"os"
)

var (
	_ logger.Logger = (*LogxAdaptor)(nil)
)

type (
	LogxAdaptor struct {
		logx.RawLogger
	}
)

func NewLogxAdaptor(name string) *LogxAdaptor {
	return &LogxAdaptor{
		RawLogger: logx.GetRawLogger(name),
	}
}

func (l *LogxAdaptor) Info(args ...interface{}) {
	l.RawLogger.Info("", args...)
}

func (l *LogxAdaptor) Warn(args ...interface{}) {
	l.RawLogger.Warn("", args...)
}

func (l *LogxAdaptor) Error(args ...interface{}) {
	l.RawLogger.Error("", args...)
}

func (l *LogxAdaptor) Debug(args ...interface{}) {
	l.RawLogger.Debug("", args...)
}

func (l *LogxAdaptor) Fatal(args ...interface{}) {
	l.RawLogger.Error("", args...)
	os.Exit(1)
}

func (l *LogxAdaptor) Infof(fmt string, args ...interface{}) {
	l.RawLogger.Infof(fmt, args...)
}

func (l *LogxAdaptor) Warnf(fmt string, args ...interface{}) {
	l.RawLogger.Warnf(fmt, args...)
}

func (l *LogxAdaptor) Errorf(fmt string, args ...interface{}) {
	l.RawLogger.Errorf(fmt, args...)
}

func (l *LogxAdaptor) Debugf(fmt string, args ...interface{}) {
	l.RawLogger.Debugf(fmt, args...)
}

func (l *LogxAdaptor) Fatalf(fmt string, args ...interface{}) {
	l.RawLogger.Errorf(fmt, args...)
	os.Exit(1)
}

func init() {
	const (
		loggerDriver = "logx"
		loggerName   = "dubbox"
	)

	extension.SetLogger(loggerDriver, func(_ *common.URL) (logger.Logger, error) {
		return NewLogxAdaptor(loggerName), nil
	})
}
