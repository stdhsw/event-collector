// Package app provides the core application logic of the event-collector.
// It assembles Exporters and the Kubernetes client into a Collector,
// runs the event collection loop, and handles graceful shutdown on OS signals.
package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/stdhsw/event-collector/internal/config"
	"github.com/stdhsw/event-collector/internal/exporter"
	"github.com/stdhsw/event-collector/internal/exporter/elasticsearch"
	"github.com/stdhsw/event-collector/internal/exporter/kafka"
	"github.com/stdhsw/event-collector/internal/exporter/volume"
	"github.com/stdhsw/event-collector/internal/kube"
	"github.com/stdhsw/event-collector/internal/logger"
)

const (
	DefaultTimeout = 10 * time.Second // exporter 종료 대기 최대 시간
)

// Collector kubernetes event를 수집하여 exporter로 전달하는 핵심 컴포넌트다.
type Collector struct {
	k8sClient *kube.Client
	exporters []exporter.Exporter
}

// NewCollector c의 설정으로 Collector를 생성하여 반환한다.
// exporter 또는 kubernetes client 생성에 실패하면 error를 반환한다.
func NewCollector(c *config.Config) (*Collector, error) {
	exporters, err := createExporters(c)
	if err != nil {
		return nil, err
	}

	// Handler는 kubernetes informer 콜백을 수신하여 exporter로 이벤트를 전달한다
	handler := &Handler{exporters: exporters}

	client, err := kube.NewClient(handler,
		kube.WithKubeConfig(c.Kube.Config),
		kube.WithResycPeriod(c.Kube.Resync),
		kube.WithNamespaces(c.Kube.Namespaces),
	)
	if err != nil {
		return nil, err
	}

	return &Collector{
		k8sClient: client,
		exporters: exporters,
	}, nil
}

// Run Collector를 실행한다. SIGINT 또는 SIGTERM 수신 시 graceful shutdown을 수행한다.
func (c *Collector) Run() {
	logger.Info("kubernetes event collector started ...")
	defer logger.Info("kubernetes event collector stopped ...")

	// context: exporter goroutine 종료 신호 전파에 사용
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 함수 종료 시 context를 반드시 해제하여 goroutine 누수 방지

	// OS signal 수신 채널 등록 및 해제
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan) // 더 이상 signal을 수신하지 않도록 정리

	var wg sync.WaitGroup

	// 각 exporter를 별도 goroutine으로 실행
	// wg.Add(1)은 goroutine 시작 전에 호출하여 race condition을 방지한다.
	for _, e := range c.exporters {
		wg.Add(1)
		go e.Start(ctx, &wg)
	}

	// kubernetes informer 시작 (blocking)
	c.k8sClient.Run()

	// SIGINT / SIGTERM 대기
	<-sigChan
	logger.Info("shutdown signal received")

	// kubernetes informer 중단 후 context 취소로 exporter goroutine 종료 신호 전달
	c.k8sClient.Close()
	cancel()

	// 모든 exporter가 종료될 때까지 최대 DefaultTimeout 대기
	exitChan := make(chan struct{})
	go func() {
		defer close(exitChan)
		wg.Wait()
	}()

	select {
	case <-exitChan:
		logger.Info("all exporters stopped")
	case <-time.After(DefaultTimeout):
		logger.Warn("timeout waiting for exporters to stop")
	}
}

// createExporters c의 설정에 따라 활성화된 exporter 목록을 생성하여 반환한다.
// 하나라도 생성에 실패하면 error를 반환한다.
func createExporters(c *config.Config) ([]exporter.Exporter, error) {
	// 최대 3개(kafka, elasticsearch, volume)를 수용할 수 있도록 capacity를 미리 할당
	exporters := make([]exporter.Exporter, 0, 3)

	if c.Kafka.Enable {
		kafkaExporter, err := kafka.NewKafkaExporter(
			c.Kafka.Brokers, c.Kafka.Topic,
			kafka.WithTimeout(c.Kafka.Timeout),
			kafka.WithRetry(c.Kafka.Retry),
			kafka.WithRetryBackoff(c.Kafka.RetryBackoff),
			kafka.WithFlushMsg(c.Kafka.FlushMsg),
			kafka.WithFlushTime(c.Kafka.FlushTime),
			kafka.WithFlushByte(c.Kafka.FlushByte),
		)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, kafkaExporter)
	}

	if c.ElasticSearch.Enable {
		elasticExporter, err := elasticsearch.NewElasticsearchExporter(
			c.ElasticSearch.Addresses, c.ElasticSearch.Index,
			elasticsearch.WithChanSize(c.ElasticSearch.ChanSize),
			elasticsearch.WithFlushTime(c.ElasticSearch.FlushTime),
			elasticsearch.WithFlushSize(c.ElasticSearch.FlushSize),
		)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, elasticExporter)
	}

	if c.Volume.Enable {
		volumeExporter, err := volume.NewVolumeExporter(
			c.Volume.FileName, c.Volume.FilePath,
			volume.WithMaxFileSize(c.Volume.MaxFileSize),
			volume.WithMaxFileCount(c.Volume.MaxFileCount),
		)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, volumeExporter)
	}

	return exporters, nil
}
