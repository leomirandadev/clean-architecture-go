package tracer

import (
	"context"
)

type Provider interface {
	Shutdown(ctx context.Context) error
	GetServiceName() string
	GetName() string
}

type Tracer struct {
	provider Provider
}

type ITracer interface {
	Close()
	GetProviderName() string
}
