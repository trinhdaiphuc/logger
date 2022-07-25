package main

import (
	"context"
	"fmt"
	"github.com/trinhdaiphuc/logger"
	pb "github.com/trinhdaiphuc/logger/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type HelloService struct{}

var _ pb.HelloServiceServer = (*HelloService)(nil)

func (h HelloService) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	log := logger.GetLogger(ctx)
	log.AddLog("Hello service")
	if len(request.Name) == 0 {
		log.AddLog("empty name")
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}
	return &pb.HelloResponse{Message: "Hello " + request.Name}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(
		logger.GrpcInterceptor(logger.ConfigGrpc{
			SkipperGrpc: func(ctx context.Context, info *grpc.UnaryServerInfo) bool {
				fmt.Println("method", info.FullMethod)
				if strings.HasSuffix(info.FullMethod, "/Hello") {
					return true
				}
				return false
			},
		}),
	))

	pb.RegisterHelloServiceServer(server, &HelloService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	// Block main routine until a signal is received
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-c

	fmt.Println("Gracefully shutting down...")
	server.GracefulStop()
	listener.Close()
	fmt.Println("GRPC was successful shutdown.")
}
