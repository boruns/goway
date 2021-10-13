package main

import (
	"fmt"
	"goway/proxy/grpc_interceptor"
	"goway/proxy/loadbalance"
	proxy2 "goway/proxy/proxy"
	"goway/proxy/public"
	"log"
	"net"
	"time"

	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

var addr = ":50051"

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	rb := loadbalance.LoadBalanceFactory(loadbalance.LbRandom)
	rb.Add("127.0.0.1:50055")

	counter, _ := public.NewFlowCountService("local_app", time.Second)
	grpcHandler := proxy2.NewGrpcLoadBalanceHandler(rb)

	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpc_interceptor.GrpcAuthStreamInterceptor,
			grpc_interceptor.GrpcFlowCountStreamInterceptor(counter),
		),
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(grpcHandler),
	)
	fmt.Println("listen at: ", addr)
	s.Serve(lis)
}
