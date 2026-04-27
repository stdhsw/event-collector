// Package app implements the Kubernetes ResourceEventHandler for the event-collector.
// Handler receives OnAdd, OnUpdate, and OnDelete callbacks from the Informer
// and forwards the converted event payload to all registered Exporters.
package app

import (
	"fmt"

	"github.com/stdhsw/event-collector/internal/exporter"
	"github.com/stdhsw/event-collector/internal/kube"
	"github.com/stdhsw/event-collector/internal/logger"
	"go.uber.org/zap"
	v1 "k8s.io/api/events/v1"
	"k8s.io/client-go/tools/cache"
)

// Handler kubernetes informer의 이벤트 콜백을 수신하여 exporter로 전달하는 핸들러다.
type Handler struct {
	exporters []exporter.Exporter
}

// OnAdd obj로 전달된 kubernetes event를 변환하여 모든 exporter에 전달한다.
func (h *Handler) OnAdd(obj any, _ bool) {
	h.handle("OnAdd", obj)
}

// OnUpdate newObj로 전달된 kubernetes event를 변환하여 모든 exporter에 전달한다.
func (h *Handler) OnUpdate(_, newObj any) {
	h.handle("OnUpdate", newObj)
}

// OnDelete obj로 전달된 kubernetes event를 변환하여 모든 exporter에 전달한다.
func (h *Handler) OnDelete(obj any) {
	h.handle("OnDelete", obj)
}

// handle op 이름과 obj를 받아 event를 변환하고 모든 exporter에 전달하는 공통 처리 함수다.
func (h *Handler) handle(op string, obj any) {
	// informer가 객체를 만료(evict)한 경우 DeletedFinalStateUnknown으로 래핑되어 전달된다.
	// 이 경우 내부의 실제 객체를 꺼내 처리한다.
	if d, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		obj = d.Obj
	}

	object, ok := obj.(*v1.Event)
	if !ok {
		logger.Warn("unexpected object type in "+op, zap.String("type", fmt.Sprintf("%T", obj)))
		return
	}

	event, err := kube.ConvertBytes(object)
	if err != nil {
		logger.Error("failed to convert event object",
			zap.String("op", op),
			zap.Error(err),
		)
		return
	}

	logger.Debug("kubernetes event received",
		zap.String("op", op),
		zap.String("event", object.Name),
		zap.String("namespace", object.Namespace),
		zap.String("kind", object.Regarding.Kind),
	)

	for _, e := range h.exporters {
		e.Write(event)
	}
}
