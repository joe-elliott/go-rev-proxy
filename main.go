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
	proxyUrl      string
	listenAddress string
	redisAddress  string
)

func init() {
	flag.StringVar(&proxyUrl, "proxy-url", "http://localhost:8081", "The address to listen on for Prometheus scrapes.")
	flag.StringVar(&listenAddress, "listen-address", ":8080", "The address to listen on for Prometheus scrapes.")
	flag.StringVar(&redisAddress, "redis-address", "redis:6379", "The address to communicate to redis on.")
}

func main() {
	flag.Parse()

	transport := &proxy.PluggableTransport{}

	transport.AddHandler(handlers.TracingHandlerFactoryFactory("go-rev-proxy"))
	transport.AddHandler(handlers.MetricsHandlerFactory)
	transport.AddHandler(handlers.TimingHandlerFactory)
	transport.AddHandler(handlers.LoggingHandlerFactory)
	transport.AddHandler(handlers.CachingHandlerFactoryFactory(redisAddress))
	transport.BuildHandlers()

	proxy := proxy.NewReverseProxy(proxyUrl, transport)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", proxy.Handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
