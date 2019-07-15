package handlers

import (
	"bufio"
	"fmt"
	"go-rev-proxy/proxy"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
)

const (
	RedisCachePrefix = "cache:"
)

func CachingHandlerFactory(cacheAddress string) proxy.TransportHandlerFactory {

	client := redis.NewClient(&redis.Options{
		Addr:     cacheAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		return func(request *http.Request, ctx *proxy.TransportHandlerContext) (*http.Response, error) {

			if request.Method != "GET" {
				return next(request, ctx)
			}

			key := genKey(request.URL)

			span := ctx.Tracer.StartSpan("CacheGet", opentracing.ChildOf(ctx.CurrentSpan.Context()))

			val, err := client.Get(key).Result()

			span.Finish()

			if err == redis.Nil {
				fmt.Println("Caching: Miss", key)
			} else if err != nil {
				return nil, err
			} else {
				fmt.Println("Caching: Hit", key)
				resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(val)), request)

				return resp, err
			}

			resp, err := next(request, ctx)

			if err != nil {
				return nil, err
			}

			respBytes, err := httputil.DumpResponse(resp, true)

			if err != nil {
				return nil, err
			}

			span = ctx.Tracer.StartSpan("CacheSet", opentracing.ChildOf(ctx.CurrentSpan.Context()))
			client.Set(key, string(respBytes), 60*time.Second).Err()
			span.Finish()

			return resp, err
		}
	}
}

func genKey(u *url.URL) string {
	return fmt.Sprintf("%v:%v", RedisCachePrefix, u)
}
