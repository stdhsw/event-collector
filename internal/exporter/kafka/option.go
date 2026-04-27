// Package kafka defines the configuration options for the Kafka exporter.
// It wraps sarama.Config with functional options for timeout, retry, partitioner,
// compression, and flush behavior.
package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

type config struct {
	saramaCfg *sarama.Config
}

func defaultConfig() *config {
	saramaCfg := sarama.NewConfig()
	// 기본 설정 (사용자 설정 가능)
	saramaCfg.Producer.MaxMessageBytes = 1024 * 1024             // 최대 메시지 크기 1MB
	saramaCfg.Producer.RequiredAcks = sarama.WaitForLocal        // 리더 브로커만 ACK 반환
	saramaCfg.Producer.Timeout = 3 * time.Second                 // 3초 타임아웃
	saramaCfg.Producer.Retry.Max = 5                             // 최대 재시도 횟수 5회
	saramaCfg.Producer.Retry.Backoff = 100 * time.Millisecond    // 재시도 간격 100ms
	saramaCfg.Producer.Partitioner = sarama.NewRandomPartitioner // 랜덤 파티션 선택

	// 메시지 압축 설정 (사용자 설정 가능)
	saramaCfg.Producer.Compression = sarama.CompressionSnappy            // 메시지 압축
	saramaCfg.Producer.CompressionLevel = sarama.CompressionLevelDefault // 압축 레벨

	// Flush 설정 (사용자 설정 가능)
	saramaCfg.Producer.Flush.Frequency = 500 * time.Millisecond // 0.5초마다 전송
	saramaCfg.Producer.Flush.Bytes = 1024 * 1024                // 1MB마다 전송
	saramaCfg.Producer.Flush.MaxMessages = 1000                 // 1000개마다 전송

	// Transaction 설정 (사용자 설정 불가능)
	saramaCfg.Producer.Transaction.Timeout = 5 * time.Second              // 트랜잭션 타임아웃 5초
	saramaCfg.Producer.Transaction.Retry.Max = 5                          // 최대 재시도 횟수 5회
	saramaCfg.Producer.Transaction.Retry.Backoff = 100 * time.Millisecond // 재시도 간격 100ms

	// ACK 반환 설정 (사용자 설정 불가능)
	saramaCfg.Producer.Return.Successes = true // 성공 시 메시지 반환
	saramaCfg.Producer.Return.Errors = true    // 전송 실패 시 에러 반환

	return &config{
		saramaCfg: saramaCfg,
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

// Option kafka exporter 설정을 변경하는 함수 타입이다.
type Option func(*config)

// WithMaxMessageBytes 메시지 최대 크기 설정
func WithMaxMessageBytes(maxMessageBytes int) Option {
	return func(c *config) {
		if maxMessageBytes > 0 {
			c.saramaCfg.Producer.MaxMessageBytes = maxMessageBytes
		}
	}
}

// WithRequiredAcks ACK 반환 설정
// 0: NoResponse, 1: WaitForLocal, -1: WaitForAll
// default: WaitForLocal
func WithRequiredAcks(requiredAcks int16) Option {
	return func(c *config) {
		ack := sarama.RequiredAcks(requiredAcks)
		switch ack {
		case sarama.WaitForAll, sarama.WaitForLocal, sarama.NoResponse:
			c.saramaCfg.Producer.RequiredAcks = ack
		default:
			c.saramaCfg.Producer.RequiredAcks = sarama.WaitForLocal
		}
	}
}

// WithTimeout 타임아웃 설정
// timeout: 최소값 1초
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		if timeout > time.Second {
			c.saramaCfg.Producer.Timeout = timeout
		}
	}
}

// WithRetry 최대 재시도 횟수 설정
// max: 최소값 0
func WithRetry(max int) Option {
	return func(c *config) {
		if max >= 0 {
			c.saramaCfg.Producer.Retry.Max = max
		}
	}
}

// WithRetryBackoff 재시도 간격 설정
// backoff: 최소값 0
func WithRetryBackoff(backoff time.Duration) Option {
	return func(c *config) {
		if backoff >= 0 {
			c.saramaCfg.Producer.Retry.Backoff = backoff
		}
	}
}

// WithPartitioner 파티션 선택 설정
// partitioner: 0(Random), 1(RoundRobin), 2(Hash)
// default: Random
func WithPartitioner(partitioner int) Option {
	return func(c *config) {
		switch partitioner {
		case 0:
			c.saramaCfg.Producer.Partitioner = sarama.NewRandomPartitioner
		case 1:
			c.saramaCfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		case 2:
			c.saramaCfg.Producer.Partitioner = sarama.NewHashPartitioner
		default:
			c.saramaCfg.Producer.Partitioner = sarama.NewRandomPartitioner
		}
	}
}

// WithCompression 메시지 압축 설정
// compression: 0(None), 1(GZIP), 2(Snappy), 3(LZ4), 4(ZSTD)
// default: None
func WithCompression(compression int) Option {
	return func(c *config) {
		switch compression {
		case 0:
			c.saramaCfg.Producer.Compression = sarama.CompressionNone
		case 1:
			c.saramaCfg.Producer.Compression = sarama.CompressionGZIP
		case 2:
			c.saramaCfg.Producer.Compression = sarama.CompressionSnappy
		case 3:
			c.saramaCfg.Producer.Compression = sarama.CompressionLZ4
		case 4:
			c.saramaCfg.Producer.Compression = sarama.CompressionZSTD
		default:
			c.saramaCfg.Producer.Compression = sarama.CompressionNone
		}
	}
}

// WithCompressionLevel 압축 레벨 설정
// level: -1(Default), 0~9
// default: Default
func WithCompressionLevel(level int) Option {
	return func(c *config) {
		if level >= -1 {
			c.saramaCfg.Producer.CompressionLevel = int(level)
		}
	}
}

// WithFlushTime 전송 주기 설정
// frequency: 최소값 0
func WithFlushTime(frequency time.Duration) Option {
	return func(c *config) {
		if frequency >= 0 {
			c.saramaCfg.Producer.Flush.Frequency = frequency
		}
	}
}

// WithFlushByte 전송 크기 설정
// bytes: 최소값 0
func WithFlushByte(bytes int) Option {
	return func(c *config) {
		if bytes >= 0 {
			c.saramaCfg.Producer.Flush.Bytes = int(bytes)
		}
	}
}

// WithFlushMsg 전송 개수 설정
// max: 최소값 0
func WithFlushMsg(max int) Option {
	return func(c *config) {
		if max >= 0 {
			c.saramaCfg.Producer.Flush.MaxMessages = max
		}
	}
}
