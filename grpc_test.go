package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	pb "github.com/trinhdaiphuc/logger/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"strings"
	"testing"
)

func TestGrpcInterceptor(t *testing.T) {
	tests := []struct {
		name       string
		reqName    string
		res        *pb.HelloResponse
		logLevel   string
		requestUri string
		config     ConfigGrpc
		errCode    codes.Code
		errMsg     string
	}{
		{
			name:       "invalid request with empty name",
			reqName:    "",
			res:        nil,
			requestUri: "/hello.HelloService/Hello",
			config:     ConfigGrpc{},
			logLevel:   "error",
			errCode:    codes.InvalidArgument,
			errMsg:     fmt.Sprint("empty name"),
		},
		{
			name:       "valid request with non negative amount",
			reqName:    "world",
			res:        &pb.HelloResponse{Message: "Hello world"},
			logLevel:   "info",
			requestUri: "/hello.HelloService/Hello",
			config: ConfigGrpc{
				SkipperGrpc: func(ctx context.Context, info *grpc.UnaryServerInfo) bool {
					fmt.Println("method", info.FullMethod)
					if strings.HasSuffix(info.FullMethod, "/Hello") {
						return true
					}
					return false
				},
			},
			errCode: codes.OK,
			errMsg:  "",
		},
		{
			name:       "valid request with non negative amount",
			reqName:    "world",
			res:        &pb.HelloResponse{Message: "Hello world"},
			logLevel:   "info",
			requestUri: "/hello.HelloService/Hello",
			config: ConfigGrpc{
				BeforeFuncGrpc: func(ctx context.Context, info *grpc.UnaryServerInfo) {
					context.WithValue(ctx, "Method", info.FullMethod)
				},
			},
			errCode: codes.OK,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			New(WithFormatter(&JSONFormatter{}), WithOutput(buf))
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(tt.config)))
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			client := pb.NewHelloServiceClient(conn)
			request := &pb.HelloRequest{Name: tt.reqName}

			response, err := client.Hello(ctx, request)

			t.Logf("Out %v, err %v", response, err)
			t.Logf("Log output %v", buf.String())

			var data map[string]interface{}
			if len(buf.String()) > 0 {
				if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
					t.Error("unexpected error", err)
				}
				_, ok := data["STEP_1"]
				uri, uriOk := data[URIField]
				level, levelOk := data[FieldKeyLevel]

				// TEST
				assert.True(t, ok, `cannot found expected "STEP_1" field: %v`, data)
				assert.True(t, uriOk, `cannot found expected "%v" field: %v`, URIField, data)
				assert.Equal(t, tt.requestUri, uri)
				assert.True(t, levelOk, `cannot found expected "%v" field: %v`, FieldKeyLevel, data)
				assert.Equal(t, tt.logLevel, level)
			}
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

func dialer(config ConfigGrpc) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			GrpcInterceptor(config),
		),
	)

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
