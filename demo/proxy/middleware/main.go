package main

import (
	"fmt"
	"goway/proxy/middleware"
	"goway/proxy/proxy"
	"log"
	"net/http"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	reverseProxy := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003/base"
		url1, err := url.Parse(rs1)
		if err != nil {
			log.Fatal(err.Error())
		}
		rs2 := "http://127.0.0.1:2004/base"
		url2, err := url.Parse(rs2)
		if err != nil {
			log.Fatal(err.Error())
		}
		urls := []*url.URL{url1, url2}
		return proxy.NewMultipleHostReverseProxy(c, urls)
	}
	log.Println("serve starting at: " + addr)

	//初始化方法数组路由器
	sliceRouter := middleware.NewSliceRouter()

	//中间件可以充当业务逻辑代码
	sg1 := sliceRouter.Group("/base").Use(middleware.TranceLogSliceRw(), func(src *middleware.SliceRouterContext) {
		src.Rw.Write([]byte("test func"))
	})
	fmt.Printf("%#v\n", sg1)
	//请求到反向代理
	sg2 := sliceRouter.Group("/").Use(middleware.TranceLogSliceRw(), func(src *middleware.SliceRouterContext) {
		fmt.Println("reverse proxy")
		reverseProxy(src)
	})
	fmt.Printf("%#v\n", sg2)
	routerHandler := middleware.NewSliceRouterHandler(nil, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}
