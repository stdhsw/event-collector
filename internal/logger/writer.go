// Package logger exposes a global zap.Logger wrapper with convenience functions.
// It provides package-level Debug, Info, Warn, Error, Fatal, and Panic functions
// that delegate to the global logger initialized by CreateGlobalLogger.
package logger

import "go.uber.org/zap"

var (
	// writer는 전역 logger 인스턴스입니다. CreateGlobalLogger 호출 전에도 안전하게 사용할 수 있도록 기본값은 no-op logger입니다.
	writer *zap.Logger = zap.NewNop()
)

// Debug msg를 debug 레벨로 기록한다. fields로 구조화 필드를 추가할 수 있다.
func Debug(msg string, fields ...zap.Field) {
	writer.Debug(msg, fields...)
}

// Info msg를 info 레벨로 기록한다. fields로 구조화 필드를 추가할 수 있다.
func Info(msg string, fields ...zap.Field) {
	writer.Info(msg, fields...)
}

// Warn msg를 warn 레벨로 기록한다. fields로 구조화 필드를 추가할 수 있다.
func Warn(msg string, fields ...zap.Field) {
	writer.Warn(msg, fields...)
}

// Error msg를 error 레벨로 기록한다. fields로 구조화 필드를 추가할 수 있다.
func Error(msg string, fields ...zap.Field) {
	writer.Error(msg, fields...)
}

// Fatal msg를 fatal 레벨로 기록한 후 os.Exit(1)을 호출한다.
func Fatal(msg string, fields ...zap.Field) {
	writer.Fatal(msg, fields...)
}

// Panic msg를 panic 레벨로 기록한 후 panic을 발생시킨다.
func Panic(msg string, fields ...zap.Field) {
	writer.Panic(msg, fields...)
}
