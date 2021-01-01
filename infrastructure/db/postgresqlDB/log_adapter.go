package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"time"

	gormLog "gorm.io/gorm/logger"
)

// ensure we implement gormLog.Interface interface (compile error otherwise)
var _ gormLog.Interface = (*GormLogger)(nil)

const traceStr string = "[%.3fms] [rows:%v] %s"
const traceWarnStr string = "%s\n[%.3fms] [rows:%v] %s"
const traceErrStr string = "%s\n[%.3fms] [rows:%v] %s"

type GormLogger struct {
	SlowThreshold time.Duration
	Logger        log.Logger
	level         gormLog.LogLevel
}

func NewGormLogger(logger *log.Logger) *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		Logger:        *logger,
		level:         gormLog.Warn,
	}
}

func (glg *GormLogger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	glg.level = level
	return glg
}

func (glg GormLogger) Error(ctx context.Context, str string, v ...interface{}) {
	glg.logRedirect(gormLog.Error, "", v...)
}

func (glg GormLogger) Info(ctx context.Context, str string, v ...interface{}) {
	glg.logRedirect(gormLog.Info, "", v...)
}

func (glg GormLogger) Warn(ctx context.Context, str string, v ...interface{}) {
	glg.logRedirect(gormLog.Warn, "", v...)
}

func (glg GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if glg.level > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && glg.level >= gormLog.Error:
			sql, rows := fc()
			if rows == -1 {
				glg.logRedirect(gormLog.Error, traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				glg.logRedirect(gormLog.Error, traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > glg.SlowThreshold && glg.SlowThreshold != 0 && glg.level >= gormLog.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", glg.SlowThreshold)
			if rows == -1 {
				glg.logRedirect(gormLog.Warn, traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				glg.logRedirect(gormLog.Warn, traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case glg.level >= gormLog.Info:
			sql, rows := fc()
			if rows == -1 {
				glg.logRedirect(gormLog.Info, traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				glg.logRedirect(gormLog.Info, traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

func (glg GormLogger) logRedirect(level gormLog.LogLevel, msg string, v ...interface{}) {
	// idea from https://github.com/jackc/pgx/blob/master/log/logrusadapter/adapter.go

	switch level {
	case gormLog.Error:
		glg.Logger.Errorf(msg, v...)
	case gormLog.Info:
		glg.Logger.Infof(msg, v...)
	case gormLog.Warn:
		glg.Logger.Warnf(msg, v...)
	default:
		glg.Logger.Warnf(fmt.Sprintf("Unknown XORM log level: %s", msg), v)
	}
}
