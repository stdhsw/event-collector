// Package testutil provides test helpers for the event-collector.
// MakeDummy constructs a synthetic kube.Event with preset values,
// useful for unit tests and local integration testing without a live cluster.
package testutil

import (
	"time"

	"github.com/stdhsw/event-collector/internal/kube"
)

// MakeDummy reVersion을 ResourceVersion으로 설정한 테스트용 kube.Event를 생성하여 반환한다.
func MakeDummy(reVersion string) *kube.Event {
	dummy := kube.Event{}
	dummy.Metadata.Name = "test"
	dummy.Metadata.Namespace = "gh-runner"
	dummy.Metadata.UID = "test1234"
	dummy.Metadata.ResourceVersion = reVersion
	dummy.Metadata.CreationTimestamp = time.Now()

	dummy.EventTime = time.Now()
	dummy.ReportingController = "test"
	dummy.Reason = "test"

	dummy.Regarding.Kind = "test"
	dummy.Regarding.Namespace = "test"
	dummy.Regarding.Name = "test"
	dummy.Regarding.UID = "test1234"
	dummy.Regarding.ApiVersion = "test.io"
	dummy.Regarding.ResourceVersion = reVersion

	dummy.Note = "test create"
	dummy.Type = "Normal"

	dummy.DeprecatedFirstTimestamp = time.Now()
	dummy.DeprecatedLastTimestamp = time.Now()
	dummy.DeprecatedCount = 1

	return &dummy
}
