//go:build !prd
// +build !prd

// Package pprof enables the Go pprof HTTP profiling server in non-production builds (build tag: !prd).
// InitPprof starts the pprof endpoint on 0.0.0.0:6060 as a goroutine,
// allowing runtime profiling during development and testing.
package pprof

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/stdhsw/event-collector/internal/logger"
	"go.uber.org/zap"
)

// InitPprof pprof HTTP 서버를 0.0.0.0:6060에서 백그라운드 고루틴으로 시작한다.
func InitPprof() {
	// Dev 환경에서만 pprof 활성화
	go func() {
		if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
			logger.Error("[pprof] server error", zap.Error(err))
		}
	}()
}
