// Package volume defines the configuration options for the Volume exporter.
// It provides functional options for maximum file size and maximum file count
// with sensible defaults (50 MB per file, 10 files).
package volume

type config struct {
	maxFileSize  int
	maxFileCount int
	chanSize     int
}

// Option volume exporter 설정을 변경하는 함수 타입이다.
type Option func(*config)

func defaultConfig() *config {
	return &config{
		maxFileSize:  50 * 1024 * 1024, // 50MB
		maxFileCount: 10,               // 10 files
		chanSize:     200,
	}
}

// fromOptions options을 기본 설정에 순서대로 적용한 config 포인터를 반환한다.
func fromOptions(options ...Option) *config {
	c := defaultConfig()
	for _, option := range options {
		option(c)
	}

	return c
}

// WithMaxFileSize 최대 파일 크기 설정
func WithMaxFileSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.maxFileSize = size
		}
	}
}

// WithMaxFileCount 최대 파일 개수 설정
func WithMaxFileCount(count int) Option {
	return func(c *config) {
		if count > 0 {
			c.maxFileCount = count
		}
	}
}

// WithChanSize volume exporter buffer channel size 설정
func WithChanSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.chanSize = size
		}
	}
}
