package tracer

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Span(ctx context.Context, identifier string, opts ...SpanStartOption) (context.Context, trace.Span) {
	options := parseOptions(opts...)
	return otel.Tracer(identifier).Start(ctx, identifier, options...)
}

type SpanStartOption struct {
	Key   string
	Value any
}

func parseOptions(startOptions ...SpanStartOption) []trace.SpanStartOption {
	opts := make([]trace.SpanStartOption, len(startOptions))

	for i, opt := range startOptions {
		value, err := json.Marshal(opt.Value)
		if err != nil {
			continue
		}
		attr := attribute.String(opt.Key, string(value))
		opts[i] = trace.WithAttributes(attr)
	}

	return opts
}
