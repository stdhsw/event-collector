// Package kube defines the configuration options for the Kubernetes client.
// Options include kubeconfig path, informer resync period, and target namespaces.
package kube

import (
	"strings"
	"time"
)

const (
	defaultKubeConfig   = ""
	defaultResyncPeriod = 0 * time.Minute // no resync by default
	minResyncPeriod     = 10 * time.Minute
)

type config struct {
	kubeConfig   string
	resyncPeriod time.Duration
	namespaces   map[string]struct{}
}

func defaultConfig() *config {
	return &config{
		kubeConfig:   defaultKubeConfig,
		resyncPeriod: defaultResyncPeriod,
		namespaces:   make(map[string]struct{}),
	}
}

type Option func(*config)

func fromOptions(options ...Option) *config {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	return config
}

// WithKubeConfig kubeConfig 경로를 설정한다. 비어 있으면 in-cluster 설정을 사용한다.
func WithKubeConfig(kubeConfig string) Option {
	return func(c *config) {
		if kubeConfig != "" {
			c.kubeConfig = kubeConfig
		}
	}
}

// WithResycPeriod informer resync 주기를 resyncPeriod로 설정한다.
// 0이면 resync를 비활성화하고, 최솟값(10분) 미만이면 10분으로 고정한다.
func WithResycPeriod(resyncPeriod time.Duration) Option {
	return func(c *config) {
		if resyncPeriod <= 0 {
			c.resyncPeriod = 0
		} else if 0 < resyncPeriod && resyncPeriod < minResyncPeriod {
			c.resyncPeriod = minResyncPeriod
		} else {
			c.resyncPeriod = resyncPeriod
		}
	}
}

// WithNamespaces 감시할 namespace 목록을 namespaces로 설정한다.
// 비어 있으면 모든 namespace를 감시한다.
func WithNamespaces(namespaces []string) Option {
	return func(c *config) {
		if len(namespaces) == 0 {
			return
		}

		c.namespaces = make(map[string]struct{})
		for _, ns := range namespaces {
			if ns == "" {
				continue
			}

			c.namespaces[strings.ToLower(ns)] = struct{}{}
		}
	}
}