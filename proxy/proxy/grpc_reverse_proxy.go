package proxy

import (
	"context"
	"goway/proxy/loadbalance"
	"log"

	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

func NewGrpcLoadBalanceHandler(lb loadbalance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr failed")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			return ctx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
