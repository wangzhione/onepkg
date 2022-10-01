package log

import (
	"context"
	"testing"
)

// online 线上环境设置 true, preview 默认是预览环境, 不是线上环境
var online bool

// SetOnline 线上环境填入 true, preview or test or local 设置 false
func SetOnline(prod bool) {
	online = prod
}

// Perview 有些感兴趣或者重要业务信息, 并且希望不要在线上打印浪费性能
// 理想很好, 但是因为 先初始化参数, 性能一般
func Preview(ctx context.Context, format string, args ...interface{}) {
	if level <= LevelInfo && !online {
		defaultLogger.Info(ctx, format, args...)
	}
}

func Test_log(t *testing.T) {
	// log 功能打印

	ctx := SetTraceID(context.Background())

	Trace(ctx, "This is Trace %d", 1)
	Debug(ctx, "This is Debug %d", 2)
	Info(ctx, "This is Info %d", 3)
	Preview(ctx, "This is Preview %d", 4)
	Warning(ctx, "This is Warning %d", 5)
	Error(ctx, "This is Error %d", 6)
}

func Benchmark_log(b *testing.B) {
	ctx := SetTraceID(context.Background())

	SetLevel(LevelTrace)

	for i := 0; i < b.N; i++ {
		Trace(ctx, "This is Trace %d", 1+i)
		Debug(ctx, "This is Debug %d", 2+i)
		Info(ctx, "This is Info %d", 3+i)
		Preview(ctx, "This is Preview %d", 4+i)
		Warning(ctx, "This is Warning %d", 5+i)
		Error(ctx, "This is Error %d", 6+i)
	}

	// 核心数                  运行次数          纳秒时间         每次调用分片内存     分配次数
	// Benchmark_log-10    	   19566	     67988 ns/op	    2223 B/op	      35 allocs/op
}
