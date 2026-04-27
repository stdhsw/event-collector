// Package main is the entry point of the event-collector application.
// It initializes the global logger from environment variables, loads the configuration file,
// and starts the Collector.
package main

import (
	"os"
	"strconv"

	"github.com/stdhsw/event-collector/internal/app"
	"github.com/stdhsw/event-collector/internal/config"
	"github.com/stdhsw/event-collector/internal/logger"
	"github.com/stdhsw/event-collector/internal/pprof"
	"go.uber.org/zap"
)

const (
	AppName string = "event-collector" // 애플리케이션 이름 (logger 파일명에 사용)

	EnvLogLevel    string = "LOG_LEVEL"    // 최소 log 레벨 환경 변수 키 (DEBUG, INFO, WARN, ERROR 등)
	EnvLogSize     string = "LOG_SIZE"     // log 파일 최대 크기(MB) 환경 변수 키
	EnvLogAge      string = "LOG_AGE"      // log 파일 보관 기간(일) 환경 변수 키
	EnvLogBack     string = "LOG_BACK"     // 보관할 이전 log 파일 최대 개수 환경 변수 키
	EnvLogCompress string = "LOG_COMPRESS" // log 파일 gzip 압축 여부 환경 변수 키

	DefaultConfigPath string = "/etc/collector/config.yaml" // 기본 설정 파일 경로
)

// init 환경 변수에서 logger 설정을 읽어 전역 logger를 초기화하고 pprof를 시작한다.
func init() {
	logLevel := os.Getenv(EnvLogLevel)
	logSize, _ := strconv.Atoi(os.Getenv(EnvLogSize))
	logAge, _ := strconv.Atoi(os.Getenv(EnvLogAge))
	logBack, _ := strconv.Atoi(os.Getenv(EnvLogBack))
	logCompress, _ := strconv.ParseBool(os.Getenv(EnvLogCompress))

	logger.CreateGlobalLogger(
		AppName,
		logger.WithLogLevel(logLevel),
		logger.WithLogMaxSize(logSize),
		logger.WithLogMaxAge(logAge),
		logger.WithLogMaxBackups(logBack),
		logger.WithLogCompress(logCompress),
	)

	pprof.InitPprof()
}

// main DefaultConfigPath에서 설정을 로드하고 Collector를 생성하여 실행한다.
func main() {
	cfg, err := config.LoadConfig(DefaultConfigPath)
	if err != nil {
		logger.Panic("failed to load config", zap.Error(err))
	}

	collector, err := app.NewCollector(cfg)
	if err != nil {
		logger.Panic("failed to create application", zap.Error(err))
	}

	collector.Run()
}
