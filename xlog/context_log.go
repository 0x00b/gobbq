package xlog

import (
	"bytes"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type ctxLogKey struct{}

var (
	_ctxLogKey ctxLogKey
	bufPool    = &sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}
)

// func OrganizeLogMiddleware(c entity.Context, in any, ret entity.RetFunc, next entity.Handler) (out any, e error) {

// 	if _organizeLog {
// 		_, ok := c.Value(_ctxLogKey).(*logrus.Entry)
// 		if ok {
// 			next(c, in, ret)
// 			return
// 		}

// 		clog, buf := newCtxLog(c)

// 		defer func() {
// 			organizeLog(c, buf)
// 		}()

// 		c = context.WithValue(c, _ctxLogKey, clog)

// 	}

// 	next(c, in, ret)
// 	return
// }

// func organizeLog(c entity.Context, buf *bytes.Buffer) {

// 	bufStr := buf.String()
// 	if len(bufStr) == 0 {
// 		return
// 	}

// 	tlog := logrus.NewEntry(logrus.New()).WithContext(c)
// 	tlog.Logger.SetReportCaller(false)
// 	tlog.Logger.SetOutput(logrus.StandardLogger().Out)

// 	tlog.Logger.SetFormatter(organizeLogFormatter)
// 	tlog.Println("{\n" + bufStr + "}")

// 	buf.Reset()
// 	bufPool.Put(buf)
// }

// func newCtxLog(c entity.Context) (*logrus.Entry, *bytes.Buffer) {

// 	logger := logrus.New()
// 	stdLogger := logrus.StandardLogger()
// 	// logger.SetOutput(stdLogger.Out)
// 	logger.SetLevel(stdLogger.Level)
// 	logger.SetFormatter(stdLogger.Formatter)
// 	logger.SetReportCaller(stdLogger.ReportCaller)

// 	buf, _ := bufPool.Get().(*bytes.Buffer)
// 	buf.Reset()
// 	logger.SetOutput(buf)
// 	for _, hook := range AddedHooks {
// 		logger.AddHook(hook)
// 	}

// 	return logrus.NewEntry(logger).WithContext(c), buf

// }

func log(c context.Context) *logrus.Entry {
	if c != nil {
		log, ok := c.Value(_ctxLogKey).(*logrus.Entry)
		if ok {
			return log.WithContext(c)
		}
	}
	return logrus.WithContext(c)
}
