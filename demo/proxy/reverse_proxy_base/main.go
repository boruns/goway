package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	proxy_addr = "http://127.0.0.1:2003"
	port       = "2002"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	proxy, err := url.Parse(proxy_addr)
	//替换协议和主机
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host

	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r)
	if err != nil {
		log.Fatal(err)
		return
	}
	//将响应的头部原路返回
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()
	//将响应内容复制给w
	bufio.NewReader(resp.Body).WriteTo(w)
}
