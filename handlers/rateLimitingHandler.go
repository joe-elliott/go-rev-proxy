package handlers

import (
	"bytes"
	"fmt"
	"go-rev-proxy/proxy"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
)

const (
	RedisRateLimitPrefix = "limit:"
)

func RateLimitingHandlerFactory(cacheAddress string, requestsPerMinute int64) proxy.TransportHandlerFactory {

	client := redis.NewClient(&redis.Options{
		Addr:     cacheAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request, ctx *proxy.TransportHandlerContext) (*http.Response, error) {

			key := genRateLimitKey(request.Host)

			// increment and go!
			span := ctx.Tracer.StartSpan("RateLimitIncr", opentracing.ChildOf(ctx.CurrentSpan.Context()))
			pipe := client.TxPipeline()
			incr := pipe.Incr(key)
			pipe.Expire(key, time.Minute)

			_, err := pipe.Exec()
			span.Finish()

			if err != nil {
				return nil, err
			}

			count := incr.Val()

			if count > requestsPerMinute {
				/*
					jpe - not working - return 429?
				*/
				return &http.Response{
					Status:        "429",
					StatusCode:    429,
					Proto:         "HTTP/1.1",
					ProtoMajor:    1,
					ProtoMinor:    1,
					Body:          ioutil.NopCloser(bytes.NewBufferString("429 - Too Many Requests")),
					ContentLength: int64(len("429 - Too Many Requests")),
					Request:       request,
					Header:        make(http.Header, 0),
				}, nil
			}

			return next(request, ctx)
		}
	}
}

func genRateLimitKey(host string) string {
	return fmt.Sprintf("%v:%v:%v", RedisRateLimitPrefix, host, time.Now().Minute())
}
