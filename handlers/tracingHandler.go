package handlers

import (
	"go-rev-proxy/proxy"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/uber/jaeger-client-go/config"
)

func TracingHandlerFactoryFactory(serviceName string) proxy.TransportHandlerFactory {

	cfg, err := config.FromEnv()

	if err != nil {
		panic(err)
	}

	cfg.ServiceName = serviceName

	tracer, _, err := cfg.NewTracer()

	if err != nil {
		panic(err)
	}

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request) (*http.Response, error) {

			request, ht = nethttp.TraceRequest(tracer, request, nethttp.OperationName("HTTP GET: "+request.URL.String()))

			err, resp := next(request)

			ht.Finish()

			return err, resp
		}
	}
}
