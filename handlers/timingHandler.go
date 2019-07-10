package handlers

import (
	"go-rev-proxy/proxy"
	"log"
	"net/http"
	"time"
)

func TimingHandlerFactory() proxy.TransportHandlerFactory {

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request) (*http.Response, error) {
			defer timeTrack(time.Now())

			return next(request)
		}
	}
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("HTTPRequest took %v", elapsed)
}
