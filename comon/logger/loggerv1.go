package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

// v1Logger 是日志门面的另一种实现
type v1Logger struct {
	_logger *zap.Logger
}

// 全局唯一的 v1Logger 实例
var v1Log *v1Logger

// InitV1Logger 初始化 v1Logger
func InitV1Logger(logger *zap.Logger) {
	v1Log = &v1Logger{
		_logger: logger,
	}
}

// 门面函数
func InfoV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		// 如果 v1Log 未初始化，使用默认的 _logger
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Info(ctx, msg, kv...)
}

func DebugV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Debug(ctx, msg, kv...)
}

func WarnV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Warn(ctx, msg, kv...)
}

func ErrorV1(ctx context.Context, msg string, kv ...interface{}) {
	if v1Log == nil {
		v1Log = &v1Logger{
			_logger: _logger,
		}
	}
	v1Log.Error(ctx, msg, kv...)
}

func (l *v1Logger) Debug(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.DebugLevel, msg, kv...)
}

func (l *v1Logger) Info(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.InfoLevel, msg, kv...)
}

func (l *v1Logger) Warn(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.WarnLevel, msg, kv...)
}

func (l *v1Logger) Error(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.ErrorLevel, msg, kv...)
}

// kv 应该是成对的数据, 类似: name,张三,age,10,...
func (l *v1Logger) log(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	// 保证要打印的日志信息成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}

	// 从 context 中获取追踪参数
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

	// 日志行信息中增加追踪参数
	kv = append(kv, "traceid", traceId, "spanid", spanId, "pspanid", pSpanId)

	// 增加日志调用者信息, 方便查日志时定位程序位置
	funcName, file, line := getV1LoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

// getV1LoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func getV1LoggerCallerInfo() (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(3) // 回溯拿调用日志方法的业务函数的信息
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
