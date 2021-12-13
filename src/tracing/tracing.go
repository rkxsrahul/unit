package go_tracing

import (
	"fmt"
	"io"

	ot "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

// function for intialising jaeger tracing using opentracing
func InitJaeger(service string) (ot.Tracer, io.Closer) {
	config.SetConfig()
	// creating new metrics factory using prometheus
	metricsFactory := prometheus.New()
	// NewZipkinB3HTTPHeaderPropagator creates a Propagator for extracting and injecting
	// Zipkin HTTP B3 headers into SpanContexts.
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	// setting jeager configuration
	cfg := &jconfig.Configuration{
		// rpc metrics enabled
		RPCMetrics: true,
		// sampler enabled
		Sampler: &jconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		// reporter enabled
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: fmt.Sprint(config.Conf.Jaeger.Host, ":", config.Conf.Jaeger.Port),
			LogSpans:           true,
		},
	}

	// using above config creating new tracer
	tracer, closer, err := cfg.New(service, jconfig.ZipkinSharedRPCSpan(true), jconfig.Injector(ot.HTTPHeaders, propagator), jconfig.Extractor(ot.HTTPHeaders, propagator), jconfig.Logger(jaeger.StdLogger), jconfig.Metrics(metricsFactory))
	if err != nil {
		// any error during creation of new tracer
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	// passing new tracer
	return tracer, closer
}
