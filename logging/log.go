package logging

import (
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	var gitTag = "" // needs to be overridden in CI `go build -ldflags="main.gitTag=v0.1.1)"`
	log = logrus.New()
	// log.Level = logrus.DebugLevel
	if len(gitTag) > 0 {
		log.Level = logrus.InfoLevel // also set log level to Info
	}
}

// Info ...
func Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Debug ...
func Debug(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Trace(format string, v ...interface{}) {
	log.Tracef(format, v...)
}
