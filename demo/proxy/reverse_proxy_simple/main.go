package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	rs1 := "http://127.0.0.1:2003/base"
	urll, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Fatal(err1)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(urll)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
