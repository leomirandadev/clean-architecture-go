package tracer

import (
	"context"
)

func New(provider Provider) ITracer {
	return Tracer{provider: provider}
}

func (trc Tracer) Close() {
	trc.provider.Shutdown(context.Background())
}

func (trc Tracer) GetProviderName() string {
	return trc.provider.GetName()
}
