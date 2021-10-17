package logger

import (
	"context"
	"fmt"
	pb "github.com/trinhdaiphuc/logger/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

type HelloService struct{}

var _ pb.HelloServiceServer = (*HelloService)(nil)

func (h HelloService) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	logger := GetLogger(ctx)
	logger.AddLog("request %v", logger.ToJsonString(request))
	if len(request.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}
	return &pb.HelloResponse{Message: "Hello " + request.Name}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer(grpc.UnaryInterceptor(
		GrpcInterceptor,
	))

	pb.RegisterHelloServiceServer(server, &HelloService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGrpcInterceptor(t *testing.T) {
	tests := []struct {
		testName string
		reqName  string
		res      *pb.HelloResponse
		errCode  codes.Code
		errMsg   string
	}{
		{
			"invalid request with empty name",
			"",
			nil,
			codes.InvalidArgument,
			fmt.Sprint("empty name"),
		},
		{
			"valid request with non negative amount",
			"world",
			&pb.HelloResponse{Message: "Hello world"},
			codes.OK,
			"",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			request := &pb.HelloRequest{Name: tt.reqName}

			response, err := client.Hello(ctx, request)

			if response != nil {
				if response.GetMessage() != tt.res.GetMessage() {
					t.Error("response: expected", tt.res.GetMessage(), "received", response.GetMessage())
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}
