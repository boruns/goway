package main

import (
	"context"
	"flag"
	"fmt"
	pb "goway/demo/proxy/grpc_server_client/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var port = flag.Int("port", 50055, "the port to serve on")

const (
	streamingCount = 10
)

type serve struct {
}

//服务端流式
func (s *serve) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	fmt.Println("---ServerStreamingEcho---")
	fmt.Printf("request receive: %v\n", req)
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", req.Message)
		err := stream.Send(&pb.EchoResponse{Message: req.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

//客户端流式
func (s *serve) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	fmt.Printf("--- clientStreamingEcho ---\n")
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("echo last receive message")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}
		message = in.Message
		fmt.Printf("request receive: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

//双向流式
func (s *serve) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Println("--- bidrectionalStreamingEcho ---")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("request receive %v, sending echo \n", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

//一元方法
func (s *serve) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Println("--- UnaryEcho --- ")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}
	fmt.Println(md)
	fmt.Printf("request receive: %v, send echo\n", req)
	return &pb.EchoResponse{Message: req.Message}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("serve err: %v\n", err)
	}
	fmt.Printf("server listening at: %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &serve{})
	s.Serve(lis)
}
