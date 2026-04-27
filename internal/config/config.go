// Package config handles loading and validation of the event-collector configuration.
// It reads a YAML file into a Config struct and verifies that all required fields
// for the enabled exporters (Kafka, Elasticsearch, Volume) are present.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/stdhsw/event-collector/internal/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kube struct {
		Config     string        `yaml:"config"`     // 없으면 in-cluster 자동 설정
		Resync     time.Duration `yaml:"resync"`     // 없으면 재수집 안함
		Namespaces []string      `yaml:"namespaces"` // 없으면 전체 네임스페이스
	} `yaml:"kube"`

	Kafka struct {
		Enable       bool          `yaml:"enable"`  // 필수
		Brokers      []string      `yaml:"brokers"` // 필수
		Topic        string        `yaml:"topic"`   // 필수
		Timeout      time.Duration `yaml:"timeout"`
		Retry        int           `yaml:"retry"`
		RetryBackoff time.Duration `yaml:"retryBackoff"`
		FlushMsg     int           `yaml:"flushMsg"`
		FlushTime    time.Duration `yaml:"flushTime"`
		FlushByte    int           `yaml:"flushByte"`
	} `yaml:"kafka"`

	ElasticSearch struct {
		Enable    bool     `yaml:"enable"`    // 필수
		Addresses []string `yaml:"addresses"` // 필수
		Index     string   `yaml:"index"`     // 필수
		ChanSize  int      `yaml:"chanSize"`
		FlushTime int      `yaml:"flushTime"`
		FlushSize int      `yaml:"flushSize"`
	} `yaml:"elasticsearch"`

	Volume struct {
		Enable       bool   `yaml:"enable"`   // 필수
		FileName     string `yaml:"fileName"` // 필수
		FilePath     string `yaml:"filePath"` // 필수
		MaxFileSize  int    `yaml:"maxFileSize"`
		MaxFileCount int    `yaml:"maxFileCount"`
	} `yaml:"volume"`
}

// LoadConfig fileName 경로의 YAML 파일을 읽어 Config 구조체로 반환한다.
// 파일 읽기 또는 필수 항목 검증에 실패하면 error를 반환한다.
func LoadConfig(fileName string) (*Config, error) {
	config := &Config{}
	if err := readFile(fileName, config); err != nil {
		return nil, fmt.Errorf("failed config read file: %w", err)
	}
	showConfig(config)

	return config, checkConfig(config)
}

// readFile fileName 경로의 파일을 읽어 config 구조체에 언마샬한다.
func readFile(fileName string, config *Config) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return err
	}

	return nil
}

// showConfig config의 주요 설정값을 debug 레벨로 출력한다.
func showConfig(config *Config) {
	logger.Debug("kubernetes config",
		zap.String("config", config.Kube.Config),
		zap.Duration("resync", config.Kube.Resync),
	)

	if config.Kafka.Enable {
		logger.Debug("kafka exporter config",
			zap.Strings("brokers", config.Kafka.Brokers),
			zap.String("topic", config.Kafka.Topic),
			zap.Duration("timeout", config.Kafka.Timeout),
			zap.Int("retry", config.Kafka.Retry),
			zap.Duration("retryBackoff", config.Kafka.RetryBackoff),
			zap.Int("flushMsg", config.Kafka.FlushMsg),
			zap.Duration("flushTime", config.Kafka.FlushTime),
			zap.Int("flushByte", config.Kafka.FlushByte),
		)
	}

	if config.ElasticSearch.Enable {
		logger.Debug("elasticsearch exporter config",
			zap.Strings("addresses", config.ElasticSearch.Addresses),
			zap.String("index", config.ElasticSearch.Index),
		)
	}

	if config.Volume.Enable {
		logger.Debug("volume config",
			zap.String("fileName", config.Volume.FileName),
			zap.String("filePath", config.Volume.FilePath),
			zap.Int("maxFileSize", config.Volume.MaxFileSize),
			zap.Int("maxFileCount", config.Volume.MaxFileCount),
		)
	}
}

// checkConfig config의 필수 항목이 모두 설정되어 있는지 검증한다.
// 하나 이상의 exporter가 활성화되어야 하며, 활성화된 exporter의 필수 필드가 채워져 있어야 한다.
func checkConfig(config *Config) error {
	if !config.Kafka.Enable && !config.ElasticSearch.Enable && !config.Volume.Enable {
		return fmt.Errorf("at least one exporter is required")
	}

	// kafka 활성화
	if config.Kafka.Enable {
		if len(config.Kafka.Brokers) == 0 {
			return fmt.Errorf("kafka brokers is required")
		}
		if config.Kafka.Topic == "" {
			return fmt.Errorf("kafka topic is required")
		}
	}

	// ElasticSearch 활성화
	if config.ElasticSearch.Enable {
		if len(config.ElasticSearch.Addresses) == 0 {
			return fmt.Errorf("elasticsearch addresses is required")
		}
		if config.ElasticSearch.Index == "" {
			return fmt.Errorf("elasticsearch index is required")
		}
	}

	// Volume 활성화
	if config.Volume.Enable {
		if config.Volume.FileName == "" {
			return fmt.Errorf("volume fileName is required")
		}

		if config.Volume.FilePath == "" {
			return fmt.Errorf("volume filePath is required")
		}
	}

	return nil
}
