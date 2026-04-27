// Package exporter defines the interfaces for event exporters in the event-collector.
// Component provides lifecycle management via Start, and Exporter extends it
// with a Write method for delivering event data to a backend destination.
package exporter

import (
	"context"
	"sync"
)

type Component interface {
	Start(ctx context.Context, wg *sync.WaitGroup) error
}

type Exporter interface {
	Component
	Write(data []byte)
}
