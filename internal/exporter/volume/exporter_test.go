package volume

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stdhsw/event-collector/internal/logger"
	"github.com/stdhsw/event-collector/internal/testutil"
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalLogger("volume_test")

	// get current path
	path, _ := os.Getwd()
	dir := filepath.Join(path, "test")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// create volume exporter
	exporter, err := NewVolumeExporter(
		"testevent", dir,
		WithMaxFileCount(5),
		WithMaxFileSize(1024),
	)
	if err != nil {
		t.Fatal(err)
	}

	// write dummy data
	dummy := testutil.MakeDummy("1")
	dummyBytes, _ := json.Marshal(dummy)
	for i := 0; i < 10; i++ {
		if err := exporter.writeData(dummyBytes); err != nil {
			t.Fatal(err)
		}
	}
}
