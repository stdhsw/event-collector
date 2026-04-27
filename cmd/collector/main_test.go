// Package main provides entry point tests for the event-collector application.
package main

import "testing"

// TestMain 테스트 실행 전후 setup/teardown을 담당한다.
func TestMain(m *testing.M) {
	// TODO: setup

	// m.Run을 호출하여 테스트 실행
	m.Run()

	// TODO: teardown
}
