package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var addr = "127.0.0.1:2002"

func main() {
	targets := make([]*url.URL, 0)
	var err error
	rs1 := "http://127.0.0.1:2003"
	url1, err := url.Parse(rs1)
	if err != nil {
		log.Fatal(err)
		return
	}
	targets = append(targets, url1)
	rs2 := "http://127.0.0.1:2004"
	url2, err := url.Parse(rs2)
	if err != nil {
		log.Fatal(err)
		return
	}
	targets = append(targets, url2)
	proxy := NewMultipleHostsReverseProxy(targets)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	var transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, //超时时间
			KeepAlive: 30 * time.Second, //长链接超时时间
		}).DialContext,
		MaxIdleConns:          100,              //最多空闲链接数
		IdleConnTimeout:       90 * time.Second, //空闲链接超时
		TLSHandshakeTimeout:   30 * time.Second, //tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //100-continue 超时时间
	}
	//请求协调者
	director := func(req *http.Request) {
		re, _ := regexp.Compile("^/dir(.*)")
		req.URL.Path = re.ReplaceAllString(req.URL.Path, "$1")

		//随机负载均衡
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		//todo 部分章节补充1
		//todo 当对域名(非内网)反向代理时需要设置此项，当做后端反向代理时不需要
		req.Host = target.Host

		//url 地址重写：重写前：/aa 重写后/base/aa
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}
	//更改内容
	modifyFunc := func(resp *http.Response) error {
		//如果是websocket
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		var payload []byte
		var readErr error
		if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			payload, readErr = ioutil.ReadAll(gr)
			resp.Header.Del("Content-Encoding")
		} else {
			payload, readErr = ioutil.ReadAll(resp.Body)
		}
		if readErr != nil {
			return readErr
		}
		if resp.StatusCode != http.StatusOK {
			payload = []byte("StatusCode error: " + string(payload))
		}

		//因为预读了数据所以内容重新回写
		resp.Body = io.NopCloser(bytes.NewBuffer(payload))
		resp.ContentLength = int64(len(payload))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error: "+err.Error(), 500)
	}
	return &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasPrefix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
