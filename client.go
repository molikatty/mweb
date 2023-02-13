package mweb

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

type send interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	one sync.Once
	c   *http.Client
	Tr  = &http.Transport{
		MaxConnsPerHost:     5,
		IdleConnTimeout:     time.Second * 5,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   false,
	}
)

func DefaultClient() send {
	if c == nil {
		one.Do(func() {
			c = &http.Client{
				Transport: Tr,
				Timeout:   time.Second * 5,
			}
		})
	}

	return c
}

func SetProxy(py Proxy) {
	py.Proxy(Tr)
}
