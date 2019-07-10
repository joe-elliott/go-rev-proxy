package handlers

import (
	"go-rev-proxy/proxy"
	"net/http"
	"time"

	"go-rev-proxy/metrics"
)

func MetricsHandlerFactory() proxy.TransportHandlerFactory {

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request, ctx *proxy.TransportHandlerContext) (*http.Response, error) {
			start := time.Now()

			resp, err := next(request, ctx)

			elapsed := time.Since(start)

			metrics.RequestLatencyMilliseconds.WithLabelValues(request.URL.Path).Observe(float64(elapsed / time.Millisecond))

			return resp, err
		}
	}
}
