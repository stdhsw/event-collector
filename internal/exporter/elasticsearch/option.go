// Package elasticsearch defines the configuration options for the Elasticsearch exporter.
// It provides functional options for channel buffer size, flush interval, flush size,
// and credentials sourced from environment variables.
package elasticsearch

import (
	"os"
	"time"
)

const (
	EnvElasticsearchUser = "ELASTICSEARCH_USER"     // Elasticsearch 인증 사용자명 환경 변수 키
	EnvElasticsearchPass = "ELASTICSEARCH_PASSWORD" // Elasticsearch 인증 비밀번호 환경 변수 키
)

type config struct {
	user, pass string
	chanSize   int
	flushTime  time.Duration
	flushSize  int
}

type Option func(*config)

func defaultConfig() *config {
	return &config{
		user:      os.Getenv(EnvElasticsearchUser),
		pass:      os.Getenv(EnvElasticsearchPass),
		chanSize:  200,
		flushTime: 1 * time.Second, // 1초
		flushSize: 1024 * 1024,     // 1MB
	}
}

func fromOptions(opts ...Option) *config {
	c := defaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithChanSize elasticsearch exporter buffer channel size 설정
func WithChanSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.chanSize = size
		}
	}
}

// WithFlushTime elasticsearch exporter buffer flush 시간 설정
func WithFlushTime(seconds int) Option {
	return func(c *config) {
		if seconds > 0 {
			c.flushTime = time.Duration(seconds) * time.Second
		}
	}
}

// WithFlushSize elasticsearch exporter buffer flush 사이즈 설정
func WithFlushSize(byteSize int) Option {
	return func(c *config) {
		if byteSize > 0 {
			c.flushSize = byteSize
		}
	}
}
