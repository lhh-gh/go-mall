package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

// v1Logger 日志门面的另一种实现方式
// 这种方式直接在日志方法中处理 context，使用更加简洁
type v1Logger struct {
	_logger *zap.Logger
}

// 全局唯一的 v1Logger 实例
var v1Log *v1Logger

// InitV1Logger 初始化 v1Logger
// 在应用启动时调用，确保日志实例被正确初始化
func InitV1Logger(logger *zap.Logger) {
	v1Log = &v1Logger{
		_logger: logger,
	}
}

// 以下是日志门面函数，提供更简洁的调用方式
// 这些函数会自动处理 context 中的追踪信息

// InfoV1 记录信息级别的日志
func InfoV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		// 如果 v1Log 未初始化，使用默认的 _logger
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Info(ctx, msg, kv...)
}

// DebugV1 记录调试级别的日志
func DebugV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Debug(ctx, msg, kv...)
}

// WarnV1 记录警告级别的日志
func WarnV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Warn(ctx, msg, kv...)
}

// ErrorV1 记录错误级别的日志
func ErrorV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Error(ctx, msg, kv...)
}

// 以下是 v1Logger 的方法实现

// Debug 实现调试级别的日志记录
func (l *v1Logger) Debug(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.DebugLevel, msg, kv...)
}

// Info 实现信息级别的日志记录
func (l *v1Logger) Info(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.InfoLevel, msg, kv...)
}

// Warn 实现警告级别的日志记录
func (l *v1Logger) Warn(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.WarnLevel, msg, kv...)
}

// Error 实现错误级别的日志记录
func (l *v1Logger) Error(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.ErrorLevel, msg, kv...)
}

// log 是实际的日志记录方法
// 参数说明：
// - ctx: 上下文，用于获取追踪信息
// - lvl: 日志级别
// - msg: 日志消息
// - kv: 键值对形式的日志字段
func (l *v1Logger) log(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	// 确保日志字段成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}

	// 从上下文中获取追踪参数
	var traceId, spanId, pSpanId string
	if ctx != nil {
		if v := ctx.Value("traceid"); v != nil {
			traceId = v.(string)
		}
		if v := ctx.Value("spanid"); v != nil {
			spanId = v.(string)
		}
		if v := ctx.Value("pspanid"); v != nil {
			pSpanId = v.(string)
		}
	}

	// 添加追踪参数到日志字段
	kv = append(kv, "traceid", traceId, "spanid", spanId, "pspanid", pSpanId)

	// 添加调用者信息，方便定位日志来源
	funcName, file, line := getV1LoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	// 将键值对转换为 zap.Field
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

// getV1LoggerCallerInfo 获取日志调用者的信息
// 返回：函数名、文件名、行号
func getV1LoggerCallerInfo() (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(3) // 回溯3层调用栈，获取实际调用日志的代码位置
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
