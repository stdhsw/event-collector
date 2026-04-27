// Package elasticsearch implements the Exporter interface for Elasticsearch.
// ElasticsearchExporter buffers incoming events and periodically flushes them
// to an Elasticsearch index using the Bulk API for efficient batch indexing.
package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stdhsw/event-collector/internal/exporter"
	"github.com/stdhsw/event-collector/internal/logger"
	"go.uber.org/zap"
)

var _ exporter.Exporter = (*ElasticsearchExporter)(nil)

type ElasticsearchExporter struct {
	es    *elasticsearch.Client
	index string

	// buffer
	dataChan  chan []byte
	mux       *sync.Mutex
	flushTime time.Duration
	flushSize int
	buffer    []byte
}

// NewElasticsearchExporter addrs와 index로 elasticsearch exporter를 생성한다.
// opts로 추가 설정을 적용할 수 있다. 연결에 실패하면 error를 반환한다.
func NewElasticsearchExporter(addrs []string, index string, opts ...Option) (*ElasticsearchExporter, error) {
	c := fromOptions(opts...)

	esCfg := elasticsearch.Config{
		Addresses: addrs,
		Username:  c.user,
		Password:  c.pass,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TLS 인증서 검증 비활성화
			},
		},
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch exporter: %w", err)
	}

	_, err = es.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}

	e := &ElasticsearchExporter{
		es:        es,
		index:     index,
		dataChan:  make(chan []byte, c.chanSize),
		mux:       &sync.Mutex{},
		flushTime: c.flushTime,
		flushSize: c.flushSize,
		buffer:    make([]byte, 0),
	}

	return e, nil
}

// Start ctx가 취소될 때까지 이벤트를 수신하여 buffer에 쌓고 주기적으로 Elasticsearch에 flush한다.
// 종료 시 wg.Done()을 호출한다.
func (e *ElasticsearchExporter) Start(ctx context.Context, wg *sync.WaitGroup) error {
	logger.Info("[elasticsearch exporter] started")
	ticker := time.NewTicker(e.flushTime)

	defer func() {
		close(e.dataChan)
		ticker.Stop()
		e.shutdown()

		logger.Info("[elasticsearch exporter] stopped")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-e.dataChan:
			e.writeBuffer(data)
			if len(e.buffer) >= e.flushSize {
				e.writeBulk(ctx)
			}
		case <-ticker.C:
			if len(e.buffer) > 0 {
				e.writeBulk(ctx)
			}
		}
	}
}

// writeBulk buffer에 쌓인 데이터를 Elasticsearch Bulk API로 전송하고 buffer를 초기화한다.
// ctx를 전달받아 30초 timeout과 함께 Bulk 요청을 수행한다.
func (e *ElasticsearchExporter) writeBulk(ctx context.Context) {
	e.mux.Lock()
	defer e.mux.Unlock()

	bulkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	buf := bytes.NewBuffer(e.buffer)
	res, err := e.es.Bulk(buf, e.es.Bulk.WithContext(bulkCtx))
	if err != nil {
		logger.Error("[elasticsearch exporter] failed to write bulk", zap.Error(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("[elasticsearch exporter] bulk request failed", zap.String("response", res.String()))
		return
	}

	logger.Debug("[elasticsearch exporter] bulk request success")

	// buffer 초기화
	e.buffer = e.buffer[:0]
}

// writeBuffer data를 Elasticsearch Bulk 형식으로 변환하여 내부 buffer에 추가한다.
func (e *ElasticsearchExporter) writeBuffer(data []byte) {
	e.mux.Lock()
	defer e.mux.Unlock()

	result := make([]byte, 0)
	meta := fmt.Sprintf(`{"index": {"_index": "%s" } }`, e.index)
	result = append(result, meta...)
	result = append(result, '\n')
	result = append(result, data...)
	result = append(result, '\n')

	e.buffer = append(e.buffer, result...)
}

// Write data를 exporter의 수신 채널에 전달한다.
// exporter가 종료 중이면 데이터를 버린다.
func (e *ElasticsearchExporter) Write(data []byte) {
	e.dataChan <- data
}

// shutdown 남은 buffer를 flush하고 exporter를 종료한다.
// ctx가 이미 취소된 상태이므로 독립적인 context.Background()를 사용한다.
func (e *ElasticsearchExporter) shutdown() {
	e.writeBulk(context.Background())
}
