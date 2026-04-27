//go:build prd
// +build prd

// Package pprof is a no-op pprof stub for production builds (build tag: prd).
// InitPprof does nothing in production to avoid exposing profiling endpoints.
package pprof

// InitPprof 프로덕션 빌드에서는 아무 동작도 수행하지 않는다.
func InitPprof() {
	// prod에서는 아무것도 하지 않음
}
