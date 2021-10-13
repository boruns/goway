package main

import (
	"fmt"
	"goway/proxy/loadbalance"
	proxy2 "goway/proxy/proxy"
	"log"
	"net"

	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

const addr = ":50051"

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed, err: %v\n", err)
	}
	rb := loadbalance.LoadBalanceFactory(loadbalance.LbRandom)
	rb.Add("127.0.0.1:50055")

	grpcHandler := proxy2.NewGrpcLoadBalanceHandler(rb)

	s := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(grpcHandler),
	)
	fmt.Printf("start server at: %s\n", addr)
	log.Fatal(s.Serve(lis))
}
