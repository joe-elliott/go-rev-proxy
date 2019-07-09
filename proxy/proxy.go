package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	url   *url.URL
	proxy *httputil.ReverseProxy
}

func NewReverseProxy(urlString string) *ReverseProxy {
	// jpe - return err
	url, err := url.Parse(urlString)

	if err != nil {
		return nil
	}

	p := httputil.NewSingleHostReverseProxy(url)
	p.Transport = &PluggableTransport{}

	return &ReverseProxy{
		url:   url,
		proxy: p,
	}
}

func (p *ReverseProxy) Handler(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
