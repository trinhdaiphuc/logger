package main

import (
	"context"
	"fmt"
	"github.com/trinhdaiphuc/logger"
	pb "github.com/trinhdaiphuc/logger/proto/hello"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log := logger.New(logger.WithFormatter(&logger.TextFormatter{}))
	client := pb.NewHelloServiceClient(conn)

	for i := 0; i <= 5; i++ {
		req := &pb.HelloRequest{
			Name: fmt.Sprintf("user-%d", i),
		}
		resp, err := client.Hello(context.Background(), req)
		if err != nil {
			log.Error(err)
		}
		log.Info("resp", log.ToJsonString(resp))
	}
}
