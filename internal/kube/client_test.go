package kube

import (
	"testing"
	"time"

	v1 "k8s.io/api/events/v1"
)

type testHandler struct {
	t *testing.T
}

func (h *testHandler) OnAdd(obj interface{}, isInInitialList bool) {
	object := obj.(*v1.Event)
	event := ConvertEvent(object)
	h.t.Logf("OnAdd: %v", event)
}

func (h *testHandler) OnUpdate(oldObj, newObj interface{}) {
	object := newObj.(*v1.Event)
	event := ConvertEvent(object)
	h.t.Logf("OnUpdate: %v", event)
}

func (h *testHandler) OnDelete(obj interface{}) {
	object := obj.(*v1.Event)
	event := ConvertEvent(object)
	h.t.Logf("OnDelete: %v", event)
}

func TestClient(t *testing.T) {
	handler := &testHandler{
		t: t,
	}

	client, err := NewClient(
		handler,
		WithResycPeriod(0),
		// WithKubeConfig(os.Getenv("KUBECONFIG")),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.Run()
	time.Sleep(10 * time.Second)
	client.Close()
}