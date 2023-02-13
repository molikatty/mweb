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
)

type Header = http.Header

type Response struct {
	HD   Header
	Body []byte
	Code int
}

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

	respFree = sync.Pool{New: func() interface{} { return new(Response) }}

	bfFree = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
)

func Get(addr string, header Header, timeout time.Duration, run send) (*Response, error) {
	get := getFree.Get().(*http.Request)
	if err := setRequest(get, addr); err != nil {
		return nil, err
	}
	get.Header = http.Header(header)

	ctx, cannle := context.WithTimeout(context.Background(), timeout)
	defer cannle()
	get = get.WithContext(ctx)
	resp, err := run.Do(get)
	getFree.Put(get)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return handleBody(resp), nil
}

func Head(addr string, header Header, timeout time.Duration, run send) (Header, error) {
	head := headFree.Get().(*http.Request)
	if err := setRequest(head, addr); err != nil {
		return nil, err
	}
	head.Header = header

	ctx, cannle := context.WithTimeout(context.Background(), timeout)
	defer cannle()

	resp, err := run.Do(head.WithContext(ctx))
	headFree.Put(head)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Header, nil
}

func Post(addr, body string, header Header, timeout time.Duration, run send) (*Response, error) {
	post := postFree.Get().(*http.Request)
	if err := setRequest(post, addr); err != nil {
		return nil, err
	}
	post.Header = header
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

	resp, err := run.Do(post.WithContext(ctx))
	postFree.Put(post)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return handleBody(resp), nil
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

func handleBody(resp *http.Response) *Response {
	buf := bfFree.Get().(*bytes.Buffer)
	buf.Reset()
	io.Copy(buf, resp.Body)
	defer bfFree.Put(buf)
	return &Response{
		Code: resp.StatusCode,
		HD:   resp.Header,
		Body: buf.Bytes(),
	}
}
