package log

import (
	"context"

	"go.uber.org/zap"
)

func AddCallerSkip(skip int) Logger {
	return std.AddCallerSkip(skip)
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents zap from always
// reporting the wrapper code as the caller.
func (l *zapLogger) AddCallerSkip(skip int) Logger {
	lc := l.clone()
	lc.z = lc.z.WithOptions(zap.AddCallerSkip(skip))
	return lc
}

// Ctx 解析传入的 context，尝试提取关注的键值，并添加到 zap.Logger 结构化日志中.
func Ctx(ctx context.Context) Logger {
	return std.Ctx(ctx)
}

// Ctx 方法，根据 context 提取字段并添加到日志中
func (l *zapLogger) Ctx(ctx context.Context) Logger {
	lc := l.clone()

	for fieldName, extractor := range l.contextExtractors {
		if val := extractor(ctx); val != "" {
			lc.z = lc.z.With(zap.String(fieldName, val))
		}
	}

	return lc
}

func With(fields ...Field) Logger {
	return std.With(fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	lc := l.clone()
	lc.z = lc.z.With(fields...)
	return lc
}

func SetLevel(level string) {
	std.SetLevel(level)
}

func (l *zapLogger) SetLevel(level string) {
	lc := l.clone()
	lc.opts.Level = level
	if lc.z != nil {
		if parsedLevel, err := zap.ParseAtomicLevel(level); err == nil {
			lc.z = lc.z.WithOptions(zap.IncreaseLevel(parsedLevel))
		}
	}
	*l = *lc
}

// clone 深度拷贝 zapLogger.
func (l *zapLogger) clone() *zapLogger {
	copied := *l
	return &copied
}

func Debug(msg string, fields ...Field) {
	std.Debug(msg, fields...)
}
func Info(msg string, fields ...Field) {
	std.Info(msg, fields...)
}
func Warn(msg string, fields ...Field) {
	std.Warn(msg, fields...)
}
func Error(msg string, fields ...Field) {
	std.Error(msg, fields...)
}
func Panic(msg string, fields ...Field) {
	std.Panic(msg, fields...)
}
func Fatal(msg string, fields ...Field) {
	std.Fatal(msg, fields...)
}

func Debugf(format string, args ...any) {
	std.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	std.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

func Panicf(format string, args ...any) {
	std.Panicf(format, args...)
}

func Fatalf(format string, args ...any) {
	std.Fatalf(format, args...)
}

func Debugw(msg string, keyvals ...any) {
	std.Debugw(msg, keyvals...)
}

func Infow(msg string, keyvals ...any) {
	std.Infow(msg, keyvals...)
}

func Warnw(msg string, keyvals ...any) {
	std.Warnw(msg, keyvals...)
}

func Errorw(err error, msg string, keyvals ...any) {
	std.Errorw(err, msg, keyvals...)
}

func Panicw(msg string, keyvals ...any) {
	std.Panicw(msg, keyvals...)
}

func Fatalw(msg string, keyvals ...any) {
	std.Fatalw(msg, keyvals...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.z.Debug(msg, fields...)
}
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.z.Info(msg, fields...)
}
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.z.Warn(msg, fields...)

}
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.z.Error(msg, fields...)
}
func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.z.Panic(msg, fields...)
}
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.z.Fatal(msg, fields...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.z.Sugar().Debugf(format, args...)
}
func (l *zapLogger) Infof(format string, args ...any) {
	l.z.Sugar().Infof(format, args...)
}
func (l *zapLogger) Warnf(format string, args ...any) {
	l.z.Sugar().Warnf(format, args...)
}
func (l *zapLogger) Errorf(format string, args ...any) {
	l.z.Sugar().Errorf(format, args...)
}
func (l *zapLogger) Panicf(format string, args ...any) {
	l.z.Sugar().Panicf(format, args...)
}
func (l *zapLogger) Fatalf(format string, args ...any) {
	l.z.Sugar().Fatalf(format, args...)
}

func (l *zapLogger) Debugw(msg string, keyvals ...any) {
	l.z.Sugar().Debugw(msg, keyvals...)
}
func (l *zapLogger) Infow(msg string, keyvals ...any) {
	l.z.Sugar().Infow(msg, keyvals...)
}
func (l *zapLogger) Warnw(msg string, keyvals ...any) {
	l.z.Sugar().Warnw(msg, keyvals...)
}
func (l *zapLogger) Errorw(err error, msg string, keyvals ...any) {
	l.z.Sugar().Errorw(msg, append(keyvals, "err", err)...)
}
func (l *zapLogger) Panicw(msg string, keyvals ...any) {
	l.z.Sugar().Panicw(msg, keyvals...)
}
func (l *zapLogger) Fatalw(msg string, keyvals ...any) {
	l.z.Sugar().Fatalw(msg, keyvals...)
}
