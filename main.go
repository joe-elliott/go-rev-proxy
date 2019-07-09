package main

import (
	"flag"
	"log"
	"net/http"

	"go-rev-proxy/handlers"
	"go-rev-proxy/proxy"
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

	transport.AddHandler(handlers.LoggingHandlerFactory)
	transport.BuildHandlers()

	proxy := proxy.NewReverseProxy(proxyUrl, transport)

	http.HandleFunc("/", proxy.Handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
