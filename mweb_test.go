package mweb

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

const (
	n = 5e4
)

var curMem uint64

var cli = &http.Client{}

func TestHttpGet(t *testing.T) {
	var g sync.WaitGroup
	for i := 0; i < n; i++ {
		g.Add(1)
		go func() {
			demoFuncGet()
			g.Done()
		}()
	}
	g.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestHttpHead(t *testing.T) {
	var g sync.WaitGroup
	for i := 0; i < n; i++ {
		g.Add(1)
		go func() {
			demoFuncHead()
			g.Done()
		}()
	}
	g.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestHttpPost(t *testing.T) {
	var g sync.WaitGroup
	for i := 0; i < n; i++ {
		g.Add(1)
		go func() {
			demoFuncPost()
			g.Done()
		}()
	}
	g.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestMwebGet(t *testing.T) {
	// SetProxy(&HttpProxy{Addr: "http://127.0.0.1:8080",})
	var g sync.WaitGroup
	for i := 0; i < n; i++ {
		g.Add(1)
		go func() {
			header := make(Header)
			header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
			Get("http://127.0.0.1:2017", header, time.Second*3)
			g.Done()
		}()
	}
	g.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestMwebPost(t *testing.T) {
	// SetProxy(&HttpProxy{Addr: "http://127.0.0.1:8080",})
	// header := make(Header)
	// header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	var g sync.WaitGroup
	for i := 0; i < n; i++ {
		g.Add(1)
		go func() {
			header := make(Header)
			header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
			Post("http://127.0.0.1:2017", "test=test", nil, time.Second*3)
			g.Done()
		}()
	}
	g.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestMwebHead(t *testing.T) {
	// SetProxy(&HttpProxy{Addr: "http://127.0.0.1:8080",})
	var g sync.WaitGroup
	for i := 0; i < n; i++ {
		g.Add(1)
		go func() {
			header := make(Header)
			header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
			Head("http://127.0.0.1:2017", header, time.Second*3)
			g.Done()
		}()
	}
	g.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func demoFuncGet() {
	rqt, err := http.NewRequest("GET", "http://127.0.0.1:2017", nil)
	if err != nil {
		return
	}
	rqt.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	ctx, cannel := context.WithTimeout(context.Background(), time.Second*3)
	defer cannel()
	rqt = rqt.WithContext(ctx)

	rsp, err := cli.Do(rqt)
	if err != nil {
		return
	}

	defer rsp.Body.Close()
	b, _ := ioutil.ReadAll(rsp.Body)
	_ = string(b)
}

func demoFuncPost() {
	rqt, err := http.NewRequest("POST", "http://www.baidu.com", bytes.NewReader([]byte("test=test")))
	if err != nil {
		return
	}
	rqt.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	ctx, cannel := context.WithTimeout(context.Background(), time.Second*3)
	defer cannel()
	rqt = rqt.WithContext(ctx)

	rsp, err := cli.Do(rqt)
	if err != nil {
		return
	}

	defer rsp.Body.Close()
	b, _ := ioutil.ReadAll(rsp.Body)
	_ = string(b)
}

func demoFuncHead() {
	rqt, err := http.NewRequest("HEAD", "http://127.0.0.1:2017", nil)
	if err != nil {
		return
	}
	rqt.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	ctx, cannel := context.WithTimeout(context.Background(), time.Second*3)
	defer cannel()
	rqt = rqt.WithContext(ctx)

	rsp, err := cli.Do(rqt)
	if err != nil {
		return
	}

	defer rsp.Body.Close()
	_ = rsp.Header
}
