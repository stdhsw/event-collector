package elasticsearch

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stdhsw/event-collector/internal/logger"
	"github.com/stdhsw/event-collector/internal/testutil"
)

const (
	envEsAddrs = "ELASTICSEARCH_ADDR"  // Elasticsearch 주소 환경 변수 키
	envEsIndex = "ELASTICSEARCH_INDEX" // Elasticsearch index 환경 변수 키
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalLogger("es_test")

	// get config from env
	esAddr := os.Getenv(envEsAddrs)
	if esAddr == "" {
		t.Fatalf("elasticsearch env %s not set, use default", envEsAddrs)
	}

	esIndex := os.Getenv(envEsIndex)
	if esIndex == "" {
		t.Fatalf("elasticsearch env %s not set, use default", envEsAddrs)
	}

	// run exporter
	es, err := NewElasticsearchExporter(
		[]string{esAddr}, esIndex,
	)
	if err != nil {
		t.Fatal(err)
	}

	// write test data
	eventData := testutil.MakeDummy("1")
	eventBytes, _ := json.Marshal(eventData)
	es.writeBuffer(eventBytes)

	eventData = testutil.MakeDummy("2")
	eventBytes, _ = json.Marshal(eventData)
	es.writeBuffer(eventBytes)

	eventData = testutil.MakeDummy("3")
	eventBytes, _ = json.Marshal(eventData)
	es.writeBuffer(eventBytes)

	es.writeBulk(t.Context())
}
