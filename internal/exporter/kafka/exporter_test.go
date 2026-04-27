package kafka

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/stdhsw/event-collector/internal/logger"
	"github.com/stdhsw/event-collector/internal/testutil"
)

const (
	envKafkaBrokers = "KAFKA_BROKERS" // Kafka broker 주소 환경 변수 키
	envKafkaTopic   = "KAFKA_TOPIC"   // Kafka topic 환경 변수 키
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalLogger("kafka_test")

	// get config from env
	kafkaBrokers := os.Getenv(envKafkaBrokers)
	if kafkaBrokers == "" {
		t.Fatalf("env %s not set, use default", envKafkaBrokers)
	}

	kafkaTopic := os.Getenv(envKafkaTopic)
	if kafkaTopic == "" {
		t.Fatalf("env %s not set, use default", envKafkaTopic)
	}

	// run exporter
	kafka, err := NewKafkaExporter([]string{kafkaBrokers}, kafkaTopic)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	go kafka.Start(ctx, &wg)

	// write test data
	eventData := testutil.MakeDummy("1")
	eventBytes, _ := json.Marshal(eventData)
	kafka.Write(eventBytes)

	eventData = testutil.MakeDummy("2")
	eventBytes, _ = json.Marshal(eventData)
	kafka.Write(eventBytes)

	// stop exporter
	cancel()
}
