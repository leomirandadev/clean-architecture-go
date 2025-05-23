package otel_jaeger

import (
	"context"
	"fmt"
	"log"

	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	jaegerLib "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semanticConvention "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Options struct {
	ServiceName string `json:"service_name" mapstructure:"service" validate:"required"`
	EndpointURL string `json:"endpoint_url" mapstructure:"endpoint" validate:"required"`
	Username    string `json:"username" mapstructure:"username"`
	Password    string `json:"password" mapstructure:"password"`
	// Sample rate is expressed as 1/X where x is the population size.
	RateSampling int `json:"rate_sampling" mapstructure:"sampling"`
}

func NewCollector(opts Options) tracer.Provider {
	exporter, err := jaegerLib.New(jaegerLib.WithCollectorEndpoint(
		jaegerLib.WithEndpoint(opts.EndpointURL),
	))

	if err != nil {
		log.Fatal("we can't initialize jaeger tracer", err)
	}

	tracerOptions := []traceSDK.TracerProviderOption{
		traceSDK.WithBatcher(exporter),
		traceSDK.WithResource(resource.NewWithAttributes(
			semanticConvention.SchemaURL,
			attribute.String("service.name", opts.ServiceName),
			attribute.String("library.language", "go"),
		)),
	}

	if opts.RateSampling > 1 {
		fractionOfTraffic := 1 / float64(opts.RateSampling)
		percentageTraffic := fractionOfTraffic * 100

		fmt.Printf("Sampling  %f percentage of traffic", percentageTraffic)

		tracerOptions = append(tracerOptions, traceSDK.WithSampler(traceSDK.TraceIDRatioBased(fractionOfTraffic)))
	}

	tracerProvider := traceSDK.NewTracerProvider(tracerOptions...)

	otel.SetTracerProvider(tracerProvider)
	return newProvider(tracerProvider, opts.ServiceName)
}

func newProvider(provider *traceSDK.TracerProvider, serviceName string) tracer.Provider {
	return &jaegerImpl{provider: provider, serviceName: serviceName}
}

type jaegerImpl struct {
	provider    *traceSDK.TracerProvider
	serviceName string
}

func (tr *jaegerImpl) Shutdown(ctx context.Context) error {
	return tr.provider.Shutdown(ctx)
}

func (tr *jaegerImpl) GetName() string {
	return "jaeger"
}

func (tr *jaegerImpl) GetServiceName() string {
	return tr.serviceName
}
