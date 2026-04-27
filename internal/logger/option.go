// Package logger defines the configuration options and defaults for the zap-based logger.
// It supports customization of log level, file rotation parameters, output encoder type,
// and local time settings via functional options.
package logger

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultPath         string = "/var/log/app/" // log file 기본 저장 디렉터리
	DefaultLogExtention string = ".log"          // log file 기본 확장자
	DefaultAppName      string = "app"           // 기본 애플리케이션 이름
	DefaultMaxSize      int    = 100             // log file 최대 크기 (MB)
	DefaultMaxBackups   int    = 3               // 보관할 이전 log file 최대 개수
	DefaultMaxAge       int    = 7               // 이전 log file 보관 기간 (일)
	DefaultLocalTime    bool   = true            // 타임스탬프에 로컬 시간 사용 여부
	DefaultCompress     bool   = false           // log file gzip 압축 여부
)

// EncoderType log 출력 형식의 종류를 나타낸다.
type EncoderType string

const (
	JSONEncoder    EncoderType = "json"    // JSON 형식으로 log를 출력한다
	ConsoleEncoder EncoderType = "console" // 사람이 읽기 쉬운 콘솔 형식으로 log를 출력한다
)

type config struct {
	appName string
	encoder zapcore.Encoder
	level   zapcore.Level
	logger  lumberjack.Logger
}

// defaultOption 기본값으로 초기화된 config 포인터를 반환한다.
func defaultOption() *config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return &config{
		appName: DefaultAppName,
		encoder: zapcore.NewJSONEncoder(encoderConfig),
		level:   zapcore.InfoLevel,
		logger: lumberjack.Logger{
			Filename:   DefaultPath + DefaultAppName + DefaultLogExtention,
			MaxSize:    DefaultMaxSize,
			MaxBackups: DefaultMaxBackups,
			MaxAge:     DefaultMaxAge,
			LocalTime:  DefaultLocalTime,
			Compress:   DefaultCompress,
		},
	}
}

// Option logger 설정을 변경하는 함수 타입이다.
type Option func(*config)

// fromOptions appName과 options을 기본 설정에 순서대로 적용한 config 포인터를 반환한다.
func fromOptions(appName string, options ...Option) *config {
	c := defaultOption()
	c.appName = appName

	for _, option := range options {
		option(c)
	}

	return c
}

// WithPath log file이 저장될 디렉터리 경로를 dirPath로 설정한다.
// dirPath가 빈 문자열이면 무시한다.
func WithPath(dirPath string) Option {
	return func(c *config) {
		if dirPath != "" {
			c.logger.Filename = filepath.Join(dirPath, c.appName+DefaultLogExtention)
		}
	}
}

// WithLogMaxSize log file 최대 크기를 maxSize(MB)로 설정한다.
// maxSize가 0 이하이면 무시한다.
func WithLogMaxSize(maxSize int) Option {
	return func(c *config) {
		if maxSize > 0 {
			c.logger.MaxSize = maxSize
		}
	}
}

// WithLogMaxBackups 보관할 이전 log file의 최대 개수를 maxBackups으로 설정한다.
// maxBackups이 0 이하이면 무시한다.
func WithLogMaxBackups(maxBackups int) Option {
	return func(c *config) {
		if maxBackups > 0 {
			c.logger.MaxBackups = maxBackups
		}
	}
}

// WithLogMaxAge 이전 log file 보관 기간을 maxAge(일)로 설정한다.
// maxAge가 0 이하이면 무시한다.
func WithLogMaxAge(maxAge int) Option {
	return func(c *config) {
		if maxAge > 0 {
			c.logger.MaxAge = maxAge
		}
	}
}

// WithLogLocalTime isLocalTime이 true이면 log file 타임스탬프에 로컬 시간을, false이면 UTC를 사용한다.
func WithLogLocalTime(isLocalTime bool) Option {
	return func(c *config) {
		c.logger.LocalTime = isLocalTime
	}
}

// WithLogCompress compress가 true이면 log file을 gzip으로 압축하고, false이면 압축하지 않는다.
func WithLogCompress(compress bool) Option {
	return func(c *config) {
		c.logger.Compress = compress
	}
}

// WithLogLevel 최소 log 레벨을 level 문자열로 설정한다.
// DEBUG, INFO, WARN/WARNING, ERROR/ERR, DPANIC, PANIC, FATAL을 지원하며 대소문자를 구분하지 않는다.
// 인식할 수 없는 값은 info 레벨로 처리한다.
func WithLogLevel(level string) Option {
	return func(c *config) {
		switch strings.ToUpper(level) {
		case "DEBUG":
			c.level = zapcore.DebugLevel
		case "WARN", "WARNING":
			c.level = zapcore.WarnLevel
		case "ERROR", "ERR":
			c.level = zapcore.ErrorLevel
		case "DPANIC":
			c.level = zapcore.DPanicLevel
		case "PANIC":
			c.level = zapcore.PanicLevel
		case "FATAL":
			c.level = zapcore.FatalLevel
		default:
			c.level = zapcore.InfoLevel
		}
	}
}

// WithEncoder log 출력 인코더를 encoderType으로 설정한다.
// ConsoleEncoder 또는 JSONEncoder(기본값)를 지원한다.
func WithEncoder(encoderType EncoderType) Option {
	return func(c *config) {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		switch encoderType {
		case ConsoleEncoder:
			c.encoder = zapcore.NewConsoleEncoder(encoderConfig)
		default:
			c.encoder = zapcore.NewJSONEncoder(encoderConfig)
		}
	}
}
