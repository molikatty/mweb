package mweb

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type Proxy interface {
	Proxy(tr *http.Transport)
}

type HttpProxy struct {
	Addr string
}

type Socks5Proxy struct {
	Addr     string
	User     string
	Password string
}

var (
	ErrHttpProxy   = errors.New(`HTTP url proxy error`)
	ErrSocks5Proxy = errors.New(`SOCKS5 url proxy error`)
)

func (h HttpProxy) Proxy(tr *http.Transport) {
	u, err := url.Parse(h.Addr)
	if err != nil {
		panic(ErrHttpProxy)
	}

	tr.Proxy = http.ProxyURL(u)
}

func (s5 Socks5Proxy) Proxy(tr *http.Transport) {
	dialer, err := proxy.SOCKS5("tcp", s5.Addr,
		&proxy.Auth{
			User:     s5.User,
			Password: s5.Password,
		},
		&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		})

	if err != nil {
		panic(ErrSocks5Proxy)
	}

	tr.Dial = dialer.Dial
}
