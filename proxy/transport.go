package proxy

import (
	"net/http"
)

type PluggableTransport struct {
}

func (t *PluggableTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	response, err := http.DefaultTransport.RoundTrip(request)

	return response, err
}
