package gorm

import (
	"fmt"
	"strings"
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/xlog"
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
	// return m
	return nil
}

func (m *GormLog) Info(c entity.Context, s string, p ...any) {
	// if m.LogLevel >= logger.Info {
	// 	//xlog.Tracef(s, p...)
	// }
}
func (m *GormLog) Warn(c entity.Context, s string, p ...any) {
	// if m.LogLevel >= logger.Warn {
	// 	//xlog.Warnf(s, p...)
	// }
}
func (m *GormLog) Error(c entity.Context, s string, p ...any) {
	// if m.LogLevel >= logger.Error {
	// 	//xlog.Errorf(s, p...)
	// }
}

func (m *GormLog) Trace(c entity.Context,
	begin time.Time, fc func() (string, int64), err error) {
	if m.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && m.LogLevel >= logger.Error:
		sql, rows := fc()
		if rows == -1 {
			xlog.Errorf(c, m.traceErrStr, sql,
				err, float64(elapsed.Nanoseconds())/1e6, "-", utils.FileWithLineNum())
		} else {
			xlog.Errorf(c, m.traceErrStr, sql,
				err, float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum())
		}
	case elapsed > m.SlowThreshold && m.SlowThreshold != 0 && m.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprint("SLOW SQL >= ", m.SlowThreshold)
		if rows == -1 {
			xlog.Warnf(c, m.traceWarnStr, sql,
				slowLog, float64(elapsed.Nanoseconds())/1e6, "-", utils.FileWithLineNum())
		} else {
			xlog.Warnf(c, m.traceWarnStr, sql,
				slowLog, float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum())
		}
	case m.LogLevel >= logger.Info:
		sql, rows := fc()
		if rows == -1 {
			xlog.Tracef(c, m.traceStr, sql,
				float64(elapsed.Nanoseconds())/1e6, "-", utils.FileWithLineNum())
		} else {
			xlog.Tracef(c, m.traceStr, sql,
				float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum())
		}
	default:
	}
}
