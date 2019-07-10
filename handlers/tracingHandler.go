package handlers

import (
	"go-rev-proxy/proxy"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func TracingHandlerFactoryFactory(serviceName string) proxy.TransportHandlerFactory {

	cfg, err := config.FromEnv()

	if err != nil {
		panic(err)
	}

	cfg.Sampler = &config.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}
	cfg.ServiceName = serviceName

	tracer, _, err := cfg.NewTracer()

	if err != nil {
		panic(err)
	}

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request) (*http.Response, error) {

			span := tracer.StartSpan("HTTPRequest", opentracing.Tag{"request", request.URL.String()})
			defer span.Finish()

			return next(request)
		}
	}
}
