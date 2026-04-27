// Package logger provides a zap-based global logger with environment-aware initialization.
// When APP_ENV is set to "dev", the global writer is a no-op logger.
// Otherwise, a production logger writing to stdout and a rotating log file is created.
package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvAppEnv   = "APP_ENV"
	EnvDevValue = "dev"
)

// isDev APP_ENV 환경 변수가 "dev"이면 true를 반환한다.
func isDev() bool {
	return os.Getenv(EnvAppEnv) == EnvDevValue
}

// CreateGlobalLogger appName과 options으로 전역 logger를 초기화한다.
// APP_ENV=dev이면 no-op logger로 설정하고, 그 외에는 stdout과 파일에 동시에 기록하는 프로덕션 logger를 생성한다.
// logger 생성에 실패하면 error를 반환한다.
func CreateGlobalLogger(appName string, options ...Option) error {
	if isDev() {
		writer = zap.NewNop()
		return nil
	}

	c := fromOptions(appName, options...)
	writer = zap.New(
		zapcore.NewCore(
			c.encoder,
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(&c.logger),
			),
			c.level,
		),
	)

	if writer == nil {
		return fmt.Errorf("failed to create logger")
	}

	return nil
}

// CreateLogger appName과 options으로 새로운 logger 인스턴스를 생성하여 반환한다.
// APP_ENV=dev이면 no-op logger를 반환하고, 그 외에는 프로덕션 logger를 생성한다.
// logger 생성에 실패하면 nil과 error를 반환한다.
func CreateLogger(appName string, options ...Option) (*zap.Logger, error) {
	if isDev() {
		return zap.NewNop(), nil
	}

	c := fromOptions(appName, options...)
	l := zap.New(
		zapcore.NewCore(
			c.encoder,
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(&c.logger),
			),
			c.level,
		),
	)

	if l == nil {
		return nil, fmt.Errorf("failed to create logger")
	}

	return l, nil
}
