package main

import (
	"goway/proxy/middleware"
	"goway/proxy/proxy"
	"log"
	"net/http"
	"net/url"
)

//熔断
func main() {
	var addr = "127.0.0.1:2002"
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003"
		url1, err := url.Parse(rs1)
		if err != nil {
			log.Fatal(err.Error())
		}
		rs2 := "http://127.0.0.1:2004"
		url2, err := url.Parse(rs2)
		if err != nil {
			log.Fatal(err.Error())
		}
		urls := []*url.URL{url1, url2}
		return proxy.NewMultipleHostReverseProxy(c, urls)
	}
	router := middleware.NewSliceRouter()
	router.Group("/").Use(middleware.RateLimiter())
	handler := middleware.NewSliceRouterHandler(coreFunc, router)
	log.Println("serve listen at: ", addr)
	log.Fatal(http.ListenAndServe(addr, handler))

}
