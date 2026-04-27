package logger

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func Test_logger(t *testing.T) {
	path, _ := os.Getwd()
	dir := filepath.Join(path, "test")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Error(err)
	}

	if err := CreateGlobalLogger(
		"test_app",
		WithPath(dir),
		WithLogLevel("debug"),
		WithLogLocalTime(false),
		WithLogCompress(true),
		WithEncoder(ConsoleEncoder),
	); err != nil {
		t.Error(err)
	}

	writer.Info("info message")
	writer.Debug("debug message", zap.String("key", "value"))
	writer.Warn("warn message", zap.String("key", "value"))
	writer.Error("error message", zap.String("key", "value"))
}