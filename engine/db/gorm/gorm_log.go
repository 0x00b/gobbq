package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0x00b/gobbq/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type LogConfig struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      logger.LogLevel
}

type GormLog struct {
	*LogConfig

	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormLog(c LogConfig) *GormLog {
	m := &GormLog{LogConfig: &c}

	m.infoStr = "%s"
	m.warnStr = m.infoStr //"%s"
	m.errStr = m.infoStr  //"%s"
	m.traceStr = "%s [%.3fms] [rows:%v] %s"
	m.traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	m.traceErrStr = m.traceWarnStr //"%s %s [%.3fms] [rows:%v] %s"

	if c.Colorful {
		m.infoStr = logger.Green + "%s " + logger.Reset + logger.Green + logger.Reset
		m.warnStr = logger.BlueBold + "%s " + logger.Reset + logger.Magenta + logger.Reset
		m.errStr = logger.Magenta + "%s " + logger.Reset + logger.Red + logger.Reset
		m.traceStr = logger.Green + "%s " + logger.Reset + logger.Yellow + "[%.3fms] " +
			logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		m.traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s " + logger.Reset + logger.RedBold +
			"[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		m.traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s " + logger.Reset +
			logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}
	return m
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseGormLevel(lvl string) logger.LogLevel {
	switch strings.ToLower(lvl) {
	case "panic":
		return logger.Error
	case "fatal":
		return logger.Error
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Warn
	case "debug":
		return logger.Info
	case "trace":
		return logger.Info
	default:
		panic("unsupport grom log level")
	}
}

func (m *GormLog) LogMode(logger.LogLevel) logger.Interface {
	return m
}

func (m *GormLog) Info(c context.Context, s string, p ...interface{}) {
	// if m.LogLevel >= logger.Info {
	// 	//log.Infof(s, p...)
	// }
}
func (m *GormLog) Warn(c context.Context, s string, p ...interface{}) {
	// if m.LogLevel >= logger.Warn {
	// 	//log.Warnf(s, p...)
	// }
}
func (m *GormLog) Error(c context.Context, s string, p ...interface{}) {
	// if m.LogLevel >= logger.Error {
	// 	//log.Errorf(s, p...)
	// }
}

func (m *GormLog) Trace(c context.Context,
	begin time.Time, fc func() (string, int64), err error) {
	if m.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && m.LogLevel >= logger.Error:
		sql, rows := fc()
		if rows == -1 {
			log.Errorf(c, m.traceErrStr, sql,
				err, float64(elapsed.Nanoseconds())/1e6, "-", utils.FileWithLineNum())
		} else {
			log.Errorf(c, m.traceErrStr, sql,
				err, float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum())
		}
	case elapsed > m.SlowThreshold && m.SlowThreshold != 0 && m.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprint("SLOW SQL >= ", m.SlowThreshold)
		if rows == -1 {
			log.Warnf(c, m.traceWarnStr, sql,
				slowLog, float64(elapsed.Nanoseconds())/1e6, "-", utils.FileWithLineNum())
		} else {
			log.Warnf(c, m.traceWarnStr, sql,
				slowLog, float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum())
		}
	case m.LogLevel >= logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Infof(c, m.traceStr, sql,
				float64(elapsed.Nanoseconds())/1e6, "-", utils.FileWithLineNum())
		} else {
			log.Infof(c, m.traceStr, sql,
				float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum())
		}
	default:
	}
}
