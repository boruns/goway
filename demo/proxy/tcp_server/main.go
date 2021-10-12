package main

import (
	"fmt"
	"goway/proxy/loadbalance"
	"goway/proxy/proxy"
	"goway/proxy/tcp_middleware"
	"goway/proxy/tcp_proxy"
)

func main() {
	//代理测试
	// rb := loadbalance.LoadBalanceFactory(loadbalance.LbRoundRobin)
	// rb.Add("127.0.0.1:2002")
	// proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	// tcpServ := tcp_proxy.TcpServer{Addr: add, Handler: proxy}
	// fmt.Println("Starting tcp_proxy ad " + add)
	// tcpServ.ListenAndServe()

	//代理redis
	rb := loadbalance.LoadBalanceFactory(loadbalance.LbRandom)
	rb.Add("127.0.0.1:6379")
	proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy}
	fmt.Println("Start tcp proxy, addr: " + addr)
	tcpServ.ListenAndServe()
}

var addr = ":7002"

// type tcpHandler struct {
// }

// func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
// 	src.Write([]byte("hello world \n"))
// 	fmt.Println(ctx.Value(tcp_proxy.LocalAddrContextKey))
// }
