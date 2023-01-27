package log_test

import (
	"context"
	"testing"

	"github.com/0x00b/gobbq/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// example for file log, using lumberjack

func TestFileLogger(t *testing.T) {

	log.Init("trace", true, true, &lumberjack.Logger{
		Filename:  "./test.log",
		MaxAge:    7,
		LocalTime: true,
	}, "", log.DefaultLogTag{})

	log.Infoln(context.Background(), "test")
}
