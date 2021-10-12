package main

import (
	"context"
	"fmt"
	"goway/proxy/tcp_proxy"
	"log"
	"net"
)

func main() {
	log.Println("start server at " + addr)
	tcpServe := tcp_proxy.TcpServer{Addr: addr, Handler: &tcpHandler{}}
	log.Fatalln(tcpServe.ListenAndServe())
}

var (
	addr = ":2002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
	fmt.Println(ctx.Value(tcp_proxy.LocalAddrContextKey))
}
