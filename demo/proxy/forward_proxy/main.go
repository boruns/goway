package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct {
}

func main() {
	fmt.Println("Serve on:8180")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("0.0.0.0:8180", nil)
}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received requrest: %s, %s, %s, \n", req.Method, req.Host, req.RemoteAddr)
	transport := http.DefaultTransport
	//浅拷贝对象
	outReq := new(http.Request)
	*outReq = *req
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	//请求下游
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	//把下游内容返回给上游
	//将头部返回
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	//将内容返回
	io.Copy(rw, res.Body)
	res.Body.Close()
}
