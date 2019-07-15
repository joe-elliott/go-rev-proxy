package main

import (
	"flag"
	"log"
	"net/http"

	"go-rev-proxy/handlers"
	"go-rev-proxy/proxy"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	proxyUrl          string
	listenAddress     string
	redisAddress      string
	requestsPerMinute int
)

func init() {
	flag.StringVar(&proxyUrl, "proxy-url", "http://localhost:8081", "The address to listen on for Prometheus scrapes.")
	flag.StringVar(&listenAddress, "listen-address", ":8080", "The address to listen on for Prometheus scrapes.")
	flag.StringVar(&redisAddress, "redis-address", "redis:6379", "The address to communicate to redis on.")
	flag.IntVar(&requestsPerMinute, "requests-per-minute", 5, "The maximum requests per minute per host.")
}

func main() {
	flag.Parse()

	transport := &proxy.PluggableTransport{}

	transport.AddHandler(handlers.TracingHandlerFactory("go-rev-proxy"))
	transport.AddHandler(handlers.MetricsHandlerFactory())
	transport.AddHandler(handlers.TimingHandlerFactory())
	transport.AddHandler(handlers.LoggingHandlerFactory())
	transport.AddHandler(handlers.RateLimitingHandlerFactory(redisAddress, int64(requestsPerMinute)))
	transport.AddHandler(handlers.CachingHandlerFactory(redisAddress))
	transport.BuildHandlers()

	proxy := proxy.NewReverseProxy(proxyUrl, transport)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", proxy.Handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
