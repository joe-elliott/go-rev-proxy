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
)

func CachingHandlerFactoryFactory(cacheAddress string) proxy.TransportHandlerFactory {

	return func(next proxy.TransportHandler) proxy.TransportHandler {

		client := redis.NewClient(&redis.Options{
			Addr:     cacheAddress,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		return func(request *http.Request) (*http.Response, error) {

			key := genKey(request.URL)
			val, err := client.Get(key).Result()

			if err == redis.Nil {
				fmt.Println("Caching: Miss", key)
			} else if err != nil {
				return nil, err
			} else {
				fmt.Println("Caching: Hit", key)
				resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(val)), request)

				return resp, err
			}

			// save in redis
			resp, err := next(request)

			if err != nil {
				return nil, err
			}

			respBytes, err := httputil.DumpResponse(resp, true)

			if err != nil {
				return nil, err
			}

			// ignore set error.  we'll just store it next time
			client.Set(key, string(respBytes), 60*time.Second).Err()

			return resp, err
		}
	}
}

func genKey(u *url.URL) string {
	return fmt.Sprintf("%v", u)
}
