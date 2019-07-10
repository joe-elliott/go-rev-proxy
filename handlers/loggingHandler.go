package handlers

import (
	"fmt"
	"go-rev-proxy/proxy"
	"net/http"
)

func LoggingHandlerFactory() proxy.TransportHandlerFactory {

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request) (*http.Response, error) {

			fmt.Printf("Starting Request %v\n", request.URL)

			resp, err := next(request)

			fmt.Printf("Ending Request %v\n", request.URL)

			return resp, err
		}
	}
}
