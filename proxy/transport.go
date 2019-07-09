package proxy

import (
	"net/http"
)

type PluggableTransport struct {
	factories   []TransportHandlerFactory
	rootHandler TransportHandler
}

type TransportHandler func(request *http.Request) (*http.Response, error)

type TransportHandlerFactory func(next TransportHandler) TransportHandler

func (t *PluggableTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	response, err := t.rootHandler(request)

	return response, err
}

func (t *PluggableTransport) AddHandler(factory TransportHandlerFactory) {
	t.factories = append(t.factories, factory)
}

func (t *PluggableTransport) BuildHandlers() {

	var currentHandler TransportHandler = finalHandler

	for i := len(t.factories) - 1; i >= 0; i-- {
		currentHandler = t.factories[i](currentHandler)
	}

	t.rootHandler = currentHandler
}

func finalHandler(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	return response, err
}