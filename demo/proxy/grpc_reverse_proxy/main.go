package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var addr = ":50051"

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("tcp listen failed: %v\n", err)
	}
	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		//拒绝一些特殊的请求
		if strings.HasPrefix(fullMethodName, "/com.example.inteernal") {
			return ctx, nil, status.Errorf(codes.Unimplemented, "Unknown method")
		}
		c, err := grpc.DialContext(ctx, "localhost:50055", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		return outCtx, c, err
	}
	s := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)
	fmt.Printf("server listen at: %s\n", addr)
	s.Serve(lis)
}
