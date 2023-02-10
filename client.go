package mweb

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

var (
	one sync.Once
	c   *http.Client
	tr  = &http.Transport{
		MaxConnsPerHost:     5,
		IdleConnTimeout:     time.Second * 5,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   false,
	}
)

func client() *http.Client {
	if c == nil {
		one.Do(func() {
			c = &http.Client{
				Transport: tr,
				Timeout:   time.Second * 5,
			}
		})
	}

	return c
}

func SetProxy(py Proxy) {
	py.Proxy(tr)
}

func SetTransport(t *http.Transport) {
	tr = t
}
