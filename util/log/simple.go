package log

import (
	"context"
	"fmt"
	"log"
)

var (
	_ Logger = (*SimpleLog)(nil)
)

// SimpleLog sample local log
type SimpleLog struct {
	logger *log.Logger
}

func (sl *SimpleLog) printf(ctx context.Context, v Level, format string, args ...interface{}) {
	// 接纳别扭
	context := fmt.Sprintf(v.String()+" "+GetTraceID(ctx)+" "+format, args...)

	sl.logger.Output(3, context)
}

// Trace 级别很低的日志级别，对于核心复用库调试会有帮助
func (sl *SimpleLog) Trace(ctx context.Context, format string, args ...interface{}) {
	sl.printf(ctx, LevelTrace, format, args...)
}

// Debug 主要用于业务开发过程中打印一些运行调试信息
func (sl *SimpleLog) Debug(ctx context.Context, format string, args ...interface{}) {
	sl.printf(ctx, LevelDebug, format, args...)
}

// Info 打印一些你感兴趣的或者重要的业务信息
func (sl *SimpleLog) Info(ctx context.Context, format string, args ...interface{}) {
	sl.printf(ctx, LevelInfo, format, args...)
}

// Warning 警告预警信息, 例如客户端参数有问题, 可能不是服务端错误, 也可能是脚本尝试
func (sl *SimpleLog) Warning(ctx context.Context, format string, args ...interface{}) {
	sl.printf(ctx, LevelWarning, format, args...)
}

// Error 直接错误信息, 开发需要重点关注
func (sl *SimpleLog) Error(ctx context.Context, format string, args ...interface{}) {
	sl.printf(ctx, LevelError, format, args...)
}
