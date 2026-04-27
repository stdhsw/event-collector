// Package volume implements the Exporter interface for local filesystem storage.
// VolumeExporter writes event data to rotating files in a specified directory,
// enforcing configurable limits on individual file size and total file count.
package volume

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stdhsw/event-collector/internal/exporter"
	"github.com/stdhsw/event-collector/internal/logger"
	"go.uber.org/zap"
)

var _ exporter.Exporter = (*VolumeExporter)(nil)

type VolumeExporter struct {
	mux         sync.Mutex
	currentFile *os.File
	fileName    string
	filePath    string

	dataChan chan []byte

	currentCount int
	maxFileSize  int
	maxFileCount int
}

// NewVolumeExporter fileName과 filePath로 volume exporter를 생성한다.
// opts로 추가 설정을 적용할 수 있다. 디렉터리 생성 또는 파일 목록 조회에 실패하면 error를 반환한다.
func NewVolumeExporter(fileName, filePath string, opts ...Option) (*VolumeExporter, error) {
	c := fromOptions(opts...)
	e := &VolumeExporter{
		fileName:     fileName,
		filePath:     filePath,
		dataChan:     make(chan []byte, c.chanSize),
		maxFileSize:  c.maxFileSize,
		maxFileCount: c.maxFileCount,
	}

	// directory 생성
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	files, err := e.getSortFileList()
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		e.currentCount = 0
	} else {
		lastCount := extractNumber(files[len(files)-1])
		e.currentCount = lastCount
	}

	return e, nil
}

// Start ctx가 취소될 때까지 이벤트를 수신하여 파일에 기록한다.
// 종료 시 wg.Done()을 호출한다.
func (e *VolumeExporter) Start(ctx context.Context, wg *sync.WaitGroup) error {
	logger.Info("[volume exporter] started")
	defer func() {
		e.shutdown()
		logger.Info("[volume exporter] stopped")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-e.dataChan:
			if err := e.writeData(data); err != nil {
				logger.Error("[volume exporter] failed to write data", zap.Error(err))
			}
		}
	}
}

// Write data를 exporter의 수신 채널에 전달한다.
func (e *VolumeExporter) Write(data []byte) {
	e.dataChan <- data
}

// shutdown 현재 파일을 닫고 채널을 종료한다.
func (e *VolumeExporter) shutdown() {
	e.mux.Lock()
	defer e.mux.Unlock()

	if e.currentFile != nil {
		e.currentFile.Close()
	}

	close(e.dataChan)
}

// getFileList filePath 디렉터리에서 fileName을 포함하는 파일 목록을 반환한다.
func (e *VolumeExporter) getFileList() ([]string, error) {
	entries, err := os.ReadDir(e.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.Contains(entry.Name(), e.fileName) {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// getSortFileList 파일 목록을 숫자 suffix 기준으로 정렬하여 반환한다.
func (e *VolumeExporter) getSortFileList() ([]string, error) {
	files, err := e.getFileList()
	if err != nil {
		return nil, err
	}

	SortByNumericSuffix(files)

	return files, nil
}

// removeFile filePath 디렉터리에서 file을 삭제한다.
func (e *VolumeExporter) removeFile(file string) error {
	if err := os.Remove(filepath.Join(e.filePath, file)); err != nil {
		return err
	}

	return nil
}

// checkAndRemove 파일 개수가 maxFileCount를 초과하면 오래된 파일부터 삭제한다.
func (e *VolumeExporter) checkAndRemove() error {
	files, err := e.getSortFileList()
	if err != nil {
		return err
	}

	if len(files) >= e.maxFileCount {
		count := len(files) - e.maxFileCount + 1
		for i := range count {
			if err := e.removeFile(files[i]); err != nil {
				logger.Error("[volume exporter] failed to remove file", zap.Error(err), zap.String("file", files[i]))
			}
		}
	}

	return nil
}

// writeData data를 현재 파일에 기록한다. 파일 크기가 maxFileSize를 초과하면 새 파일을 생성한다.
func (e *VolumeExporter) writeData(data []byte) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	if e.currentFile == nil {
		file, err := os.Create(filepath.Join(e.filePath, fmt.Sprintf("%s_%d", e.fileName, e.currentCount)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		e.currentFile = file
	}

	fInfo, err := e.currentFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	if fInfo.Size() > int64(e.maxFileSize) {
		if err := e.checkAndRemove(); err != nil {
			return fmt.Errorf("failed to check and remove: %w", err)
		}

		if err := e.currentFile.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		e.currentCount++
		file, err := os.Create(filepath.Join(e.filePath, fmt.Sprintf("%s_%d", e.fileName, e.currentCount)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		e.currentFile = file
	}

	if _, err := e.currentFile.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	logger.Debug("[volume exporter] data written", zap.String("file", e.currentFile.Name()), zap.Int("size", len(data)))

	return nil
}
