package mweb

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
	"unsafe"
)

type Header http.Header

var (
	getFree = sync.Pool{
		New: func() interface{} {
			return &http.Request{
				Method:     "GET",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Body:       nil,
			}
		},
	}

	postFree = sync.Pool{
		New: func() interface{} {
			return &http.Request{
				Method:     "POST",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
		},
	}

	headFree = sync.Pool{
		New: func() interface{} {
			return &http.Request{
				Method:     "HEAD",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
		},
	}

	bfFree = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
)

func Get(addr string, header Header, timeout time.Duration) (string, error) {
	get := getFree.Get().(*http.Request)
	if err := setRequest(get, addr); err != nil {
		return "", err
	}
	get.Header = http.Header(header)
	getFree.Put(get)

	ctx, cannle := context.WithTimeout(context.Background(), timeout)
	defer cannle()

	resp, err := client().Do(get.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := bfFree.Get().(*bytes.Buffer)
	buf.Reset()
	io.Copy(buf, resp.Body)
	bfFree.Put(buf)
	b := buf.Bytes()

	return *(*string)(unsafe.Pointer(&b)), nil

}

func Head(addr string, header Header, timeout time.Duration) (Header, error) {
	head := headFree.Get().(*http.Request)
	if err := setRequest(head, addr); err != nil {
		return nil, err
	}
	head.Header = http.Header(header)
	headFree.Put(head)

	ctx, cannle := context.WithTimeout(context.Background(), timeout)
	defer cannle()

	resp, err := client().Do(head.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return Header(resp.Header), nil
}

func Post(addr, body string, header Header, timeout time.Duration) (string, error) {
	post := postFree.Get().(*http.Request)
	if err := setRequest(post, addr); err != nil {
		return "", err
	}
	post.Header = http.Header(header)
	postFree.Put(post)
	if body == "" {
		post.Body = nil
	} else {
		by := bytes.NewReader([]byte(body))
		noClose := ioutil.NopCloser(by)
		post.Body = noClose
		post.ContentLength = int64(by.Len())
		post.GetBody = func() (io.ReadCloser, error) {
			return noClose, nil
		}
	}

	ctx, cannle := context.WithTimeout(context.Background(), timeout)
	defer cannle()

	resp, err := client().Do(post.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := bfFree.Get().(*bytes.Buffer)
	buf.Reset()
	io.Copy(buf, resp.Body)
	b := buf.Bytes()
	bfFree.Put(buf)

	return *(*string)(unsafe.Pointer(&b)), nil
}

func setRequest(h *http.Request, addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	h.URL = u
	h.Host = u.Host
	return nil
}

func (h Header) Set(key, value string) {
	h[key] = []string{value}
}
