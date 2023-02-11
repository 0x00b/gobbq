package xlog_test

import (
	"testing"

	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// example for file log, using lumberjack

func TestFileLogger(t *testing.T) {

	xlog.Init("trace", true, true, &lumberjack.Logger{
		Filename:  "./test.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	xlog.Infoln(nil, "test")
}
