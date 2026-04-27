// Package kube defines the common Event model and conversion utilities.
// ConvertEvent maps a Kubernetes v1.Event into the internal Event struct,
// and ConvertBytes serializes it to JSON bytes for downstream exporters.
package kube

import (
	"encoding/json"
	"time"

	v1 "k8s.io/api/events/v1"
)

type Event struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`

	EventTime            time.Time `json:"eventTime"`
	ReportingController string    `json:"reportingController"`
	Reason               string    `json:"reason"`

	Regarding struct {
		Kind            string `json:"kind"`
		Namespace       string `json:"namespace"`
		Name            string `json:"name"`
		UID             string `json:"uid"`
		ApiVersion      string `json:"apiVersion"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"regarding"`

	Note string `json:"note"`
	Type string `json:"type"`

	DeprecatedFirstTimestamp time.Time `json:"deprecatedFirstTimestamp"`
	DeprecatedLastTimestamp  time.Time `json:"deprecatedLastTimestamp"`
	DeprecatedCount          int       `json:"deprecatedCount"`
}

// ConvertEvent object를 내부 Event 모델로 변환하여 반환한다.
func ConvertEvent(object *v1.Event) *Event {
	event := &Event{}
	event.Metadata.Name = object.Name
	event.Metadata.Namespace = object.Namespace
	event.Metadata.UID = string(object.UID)
	event.Metadata.ResourceVersion = object.ResourceVersion
	event.Metadata.CreationTimestamp = object.CreationTimestamp.Time

	event.EventTime = object.EventTime.Time
	event.ReportingController = object.ReportingController
	event.Reason = object.Reason

	event.Regarding.Kind = object.Regarding.Kind
	event.Regarding.Namespace = object.Regarding.Namespace
	event.Regarding.Name = object.Regarding.Name
	event.Regarding.UID = string(object.Regarding.UID)
	event.Regarding.ApiVersion = object.Regarding.APIVersion
	event.Regarding.ResourceVersion = object.Regarding.ResourceVersion

	event.Note = object.Note
	event.Type = object.Type

	event.DeprecatedFirstTimestamp = object.DeprecatedFirstTimestamp.Time
	event.DeprecatedLastTimestamp = object.DeprecatedLastTimestamp.Time
	event.DeprecatedCount = int(object.DeprecatedCount)

	return event
}

// ConvertBytes object를 JSON으로 직렬화하여 반환한다.
// 직렬화에 실패하면 error를 반환한다.
func ConvertBytes(object *v1.Event) ([]byte, error) {
	event := ConvertEvent(object)
	return json.Marshal(event)
}