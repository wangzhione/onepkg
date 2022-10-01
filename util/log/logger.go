package log

import (
	"context"
	"log"
	"os"
)

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	// Trace 级别很低的日志级别，对于核心复用库调试会有帮助
	Trace(ctx context.Context, format string, args ...interface{})

	// Debug 主要用于业务开发过程中打印一些运行调试信息
	Debug(ctx context.Context, format string, args ...interface{})

	// Info 打印一些你感兴趣的或者重要的业务信息
	Info(ctx context.Context, format string, args ...interface{})

	// Warning 警告预警信息, 例如客户端参数有问题, 可能不是服务端错误, 也可能是脚本尝试
	Warning(ctx context.Context, format string, args ...interface{})

	// Error 直接错误信息, 开发需要重点关注
	Error(ctx context.Context, format string, args ...interface{})
}

// Level defines the priority of a log message.
// When a logger is configured with a level, any log message with a lower
// log level (smaller by integer comparison) will not be output.
type Level int

// The levels of logs.
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelOFF
)

// String Stringer 接口实现
func (v Level) String() string {
	switch v {
	// log system use
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelOFF:
		return "OFF"
	// It will never be executed here
	default:
		return "UNKNOW"
	}
}

var (
	// level 日志等级, 默认是 LevelTrace
	level = LevelDebug

	// defaultLogger 默认 Logger
	defaultLogger Logger = &SimpleLog{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
	}
)

// SetDefaultLogger sets the default logger.
// This is not concurrency safe, which means it should only be called during init.
func SetDefaultLogger(er Logger) {
	if er == nil {
		panic("logger must not be nil")
	}
	defaultLogger = er
}

// SetLevel sets the level of logs below which logs will not be output.
// The default log level is LevelTrace.
func SetLevel(v Level) {
	if v >= LevelTrace && v <= LevelOFF {
		level = v
	}
}

// Trace 级别很低的日志级别，对于核心复用库调试会有帮助
func Trace(ctx context.Context, format string, args ...interface{}) {
	if level <= LevelTrace {
		defaultLogger.Trace(ctx, format, args...)
	}
}

// Debug 主要用于业务开发过程中打印一些运行调试信息
func Debug(ctx context.Context, format string, args ...interface{}) {
	if level <= LevelDebug {
		defaultLogger.Debug(ctx, format, args...)
	}
}

// Info 打印一些你感兴趣的或者重要的业务信息
func Info(ctx context.Context, format string, args ...interface{}) {
	if level <= LevelInfo {
		defaultLogger.Info(ctx, format, args...)
	}
}

// Warning 警告预警信息, 例如客户端参数有问题, 可能不是服务端错误, 也可能是脚本尝试
func Warning(ctx context.Context, format string, args ...interface{}) {
	if level <= LevelWarning {
		defaultLogger.Warning(ctx, format, args...)
	}
}

// Error 直接错误信息, 开发需要重点关注
func Error(ctx context.Context, format string, args ...interface{}) {
	if level <= LevelError {
		defaultLogger.Error(ctx, format, args...)
	}
}
