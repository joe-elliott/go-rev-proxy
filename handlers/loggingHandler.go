package handlers

import (
	"fmt"
	"go-rev-proxy/proxy"
	"net/http"
)

func LoggingHandlerFactory(next proxy.TransportHandler) proxy.TransportHandler {

	return func(request *http.Request) (*http.Response, error) {

		fmt.Println("before")

		resp, err := next(request)

		fmt.Println("after")

		return resp, err
	}
}
