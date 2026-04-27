// Package kafka implements the Exporter interface for Apache Kafka.
// KafkaExporter uses the sarama AsyncProducer to deliver event data to a Kafka topic
// with configurable batching, compression, and retry settings.
package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/stdhsw/event-collector/internal/exporter"
	"github.com/stdhsw/event-collector/internal/logger"
	"go.uber.org/zap"
)

var _ exporter.Exporter = (*KafkaExporter)(nil)

type KafkaExporter struct {
	producer sarama.AsyncProducer
	topic    string
}

// NewKafkaExporter brokers와 topic으로 kafka exporter를 생성한다.
// opts로 추가 설정을 적용할 수 있다. 생성에 실패하면 error를 반환한다.
func NewKafkaExporter(brokers []string, topic string, opts ...Option) (*KafkaExporter, error) {
	c := fromOptions(opts...)

	producer, err := sarama.NewAsyncProducer(brokers, c.saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka exporter: %w", err)
	}

	return &KafkaExporter{
		producer: producer,
		topic:    topic,
	}, nil
}

// Start ctx가 취소될 때까지 Kafka producer의 성공/실패 이벤트를 처리한다.
// 종료 시 wg.Done()을 호출한다.
func (e *KafkaExporter) Start(ctx context.Context, wg *sync.WaitGroup) error {
	logger.Info("[kafka exporter] started")
	defer func() {
		e.shutdown()
		logger.Info("[kafka exporter] stopped")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-e.producer.Errors():
			logger.Error("[kafka exporter] failed to send message", zap.Error(err), zap.String("topic", err.Msg.Topic), zap.Int32("partition", err.Msg.Partition))
		case success := <-e.producer.Successes():
			logger.Debug("[kafka exporter] message sent", zap.String("topic", success.Topic), zap.Int32("partition", success.Partition), zap.Int64("offset", success.Offset))
		}
	}
}

// shutdown AsyncProducer를 닫고 exporter를 종료한다.
func (e *KafkaExporter) shutdown() {
	e.producer.AsyncClose()
}

// Write data를 Kafka topic으로 비동기 전송한다.
func (e *KafkaExporter) Write(data []byte) {
	e.producer.Input() <- &sarama.ProducerMessage{
		Topic: e.topic,
		Key:   nil,
		Value: sarama.ByteEncoder(data),
	}
}
