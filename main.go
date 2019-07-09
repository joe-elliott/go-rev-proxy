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
)

func init() {
	flag.StringVar(&proxyUrl, "proxy-url", "http://localhost:8081", "The address to listen on for Prometheus scrapes.")
	flag.StringVar(&listenAddress, "listen-address", ":8080", "The address to listen on for Prometheus scrapes.")
}

func main() {
	flag.Parse()

	transport := &proxy.PluggableTransport{}

	transport.AddHandler(handlers.TimingHandlerFactory)
	transport.AddHandler(handlers.LoggingHandlerFactory)
	transport.BuildHandlers()

	proxy := proxy.NewReverseProxy(proxyUrl, transport)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", proxy.Handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
