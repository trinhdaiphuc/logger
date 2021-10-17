# Logger

## Introduce:

A wrapper [Logrus](https://github.com/sirupsen/logrus) library which make logger attached in request context, log should
be showed as steps and in only one line for every request.

Support Logger middleware for [Echo](https://echo.labstack.com/), [Fiber](https://gofiber.io/)
, [Gin](https://github.com/gin-gonic/gin) and [GRPC](https://grpc.io/docs/languages/go/basics/) framework. You can make
your own Logger middleware by adding your logger into request context and getting it by `logger.GetLogger(ctx)`
function. You can follow my middlewares and create yours.

## Usages:

**Install the package:**

```shell
go get -u github.com/trinhdaiphuc/logger
```

### Echo

Example code:

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/logger"
)

func main() {
	server := echo.New()
	server.Use(logger.EchoMiddleware)

	server.GET("/hello/:name", func(ctx echo.Context) error {
		log := logger.GetLogger(ctx.Request().Context())
		name := ctx.Param("name")
		log.AddLog("request name %v", name)
		return ctx.String(200, "Hello "+name)
	})

	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
```

Try logger with request:

```shell
for i in {0..5}; do curl :8080/hello/user-${i}; done 
```

Logger output:

```shell
    ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.6.1
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
{"STEP_1":"request name user-0","Status":200,"client_ip":"::1","end":"2021-10-17T16:27:10.717449+07:00","level":"info","msg":"latency: 25.637µs","request_method":"GET","time":"2021-10-17T16:27:10+07:00","uri":"/hello/user-0","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-1","Status":200,"client_ip":"::1","end":"2021-10-17T16:27:10.728884+07:00","level":"info","msg":"latency: 13.214µs","request_method":"GET","time":"2021-10-17T16:27:10+07:00","uri":"/hello/user-1","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-2","Status":200,"client_ip":"::1","end":"2021-10-17T16:27:10.740955+07:00","level":"info","msg":"latency: 33.484µs","request_method":"GET","time":"2021-10-17T16:27:10+07:00","uri":"/hello/user-2","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-3","Status":200,"client_ip":"::1","end":"2021-10-17T16:27:10.752934+07:00","level":"info","msg":"latency: 23.883µs","request_method":"GET","time":"2021-10-17T16:27:10+07:00","uri":"/hello/user-3","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-4","Status":200,"client_ip":"::1","end":"2021-10-17T16:27:10.765675+07:00","level":"info","msg":"latency: 27.749µs","request_method":"GET","time":"2021-10-17T16:27:10+07:00","uri":"/hello/user-4","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-5","Status":200,"client_ip":"::1","end":"2021-10-17T16:27:10.778096+07:00","level":"info","msg":"latency: 30.309µs","request_method":"GET","time":"2021-10-17T16:27:10+07:00","uri":"/hello/user-5","user_agent":"curl/7.64.1"}
```

### Fiber

Example code:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trinhdaiphuc/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.FiberMiddleware())

	app.Get("/hello/:name", func(ctx *fiber.Ctx) error {
		log := logger.GetLogger(ctx.Context())
		name := ctx.Params("name")
		log.AddLog("request name %v", name)
		return ctx.Status(200).SendString("Hello " + name)
	})

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
```

Try logger with request:

```shell
for i in {0..5}; do curl :8080/hello/user-${i}; done 
```

Logger output:

```shell
 ┌───────────────────────────────────────────────────┐ 
 │                   Fiber v2.20.2                   │ 
 │               http://127.0.0.1:8080               │ 
 │       (bound on host 0.0.0.0 and port 8080)       │ 
 │                                                   │ 
 │ Handlers ............. 3  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 42319 │ 
 └───────────────────────────────────────────────────┘ 

{"STEP_1":"request name user-0","Status":200,"client_ip":"127.0.0.1","end":"2021-10-17T16:28:24.268446+07:00","level":"info","msg":"latency: 26.555µs","request_method":"GET","time":"2021-10-17T16:28:24+07:00","uri":"/hello/user-0","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-1","Status":200,"client_ip":"127.0.0.1","end":"2021-10-17T16:28:24.274605+07:00","level":"info","msg":"latency: 9.141µs","request_method":"GET","time":"2021-10-17T16:28:24+07:00","uri":"/hello/user-1","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-2","Status":200,"client_ip":"127.0.0.1","end":"2021-10-17T16:28:24.280196+07:00","level":"info","msg":"latency: 9.223µs","request_method":"GET","time":"2021-10-17T16:28:24+07:00","uri":"/hello/user-2","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-3","Status":200,"client_ip":"127.0.0.1","end":"2021-10-17T16:28:24.286032+07:00","level":"info","msg":"latency: 12.195µs","request_method":"GET","time":"2021-10-17T16:28:24+07:00","uri":"/hello/user-3","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-4","Status":200,"client_ip":"127.0.0.1","end":"2021-10-17T16:28:24.292232+07:00","level":"info","msg":"latency: 17.991µs","request_method":"GET","time":"2021-10-17T16:28:24+07:00","uri":"/hello/user-4","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-5","Status":200,"client_ip":"127.0.0.1","end":"2021-10-17T16:28:24.301115+07:00","level":"info","msg":"latency: 41.77µs","request_method":"GET","time":"2021-10-17T16:28:24+07:00","uri":"/hello/user-5","user_agent":"curl/7.64.1"}
```

### Gin

Example code:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trinhdaiphuc/logger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(logger.GinMiddleware())
	server.GET("/hello/:name", func(ctx *gin.Context) {
		log := logger.GetLogger(ctx)
		name := ctx.Param("name")
		log.AddLog("request name %v", name)
		ctx.String(200, "Hello "+name)
	})

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
```

Try logger with request:

```shell
for i in {0..5}; do curl :8080/hello/user-${i}; done 
```

Logger output:

```shell
{"STEP_1":"request name user-0","Status":200,"client_ip":"::1","end":"2021-10-17T16:28:52.836854+07:00","level":"info","msg":"latency: 26.574µs","request_method":"GET","time":"2021-10-17T16:28:52+07:00","uri":"/hello/user-0","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-1","Status":200,"client_ip":"::1","end":"2021-10-17T16:28:52.842935+07:00","level":"info","msg":"latency: 9.643µs","request_method":"GET","time":"2021-10-17T16:28:52+07:00","uri":"/hello/user-1","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-2","Status":200,"client_ip":"::1","end":"2021-10-17T16:28:52.848682+07:00","level":"info","msg":"latency: 14.542µs","request_method":"GET","time":"2021-10-17T16:28:52+07:00","uri":"/hello/user-2","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-3","Status":200,"client_ip":"::1","end":"2021-10-17T16:28:52.855196+07:00","level":"info","msg":"latency: 10.794µs","request_method":"GET","time":"2021-10-17T16:28:52+07:00","uri":"/hello/user-3","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-4","Status":200,"client_ip":"::1","end":"2021-10-17T16:28:52.863981+07:00","level":"info","msg":"latency: 24.512µs","request_method":"GET","time":"2021-10-17T16:28:52+07:00","uri":"/hello/user-4","user_agent":"curl/7.64.1"}
{"STEP_1":"request name user-5","Status":200,"client_ip":"::1","end":"2021-10-17T16:28:52.874219+07:00","level":"info","msg":"latency: 12.744µs","request_method":"GET","time":"2021-10-17T16:28:52+07:00","uri":"/hello/user-5","user_agent":"curl/7.64.1"}
```

### GRPC

Example code:

```go
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
		logger.GrpcInterceptor,
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
```

Try logger with request client:

```go
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
```

Logger output:

```shell
{"STEP_1":"Hello service","client_ip":"127.0.0.1:62328","code":"OK","level":"info","msg":"latency: 9.363µs","request":{"name":"user-0"},"response":{"message":"Hello user-0"},"start":"2021-10-17T16:34:56.495895+07:00","time":"2021-10-17T16:34:56+07:00","uri":"/hello.HelloService/Hello"}
{"STEP_1":"Hello service","client_ip":"127.0.0.1:62328","code":"OK","level":"info","msg":"latency: 34.926µs","request":{"name":"user-1"},"response":{"message":"Hello user-1"},"start":"2021-10-17T16:34:56.496685+07:00","time":"2021-10-17T16:34:56+07:00","uri":"/hello.HelloService/Hello"}
{"STEP_1":"Hello service","client_ip":"127.0.0.1:62328","code":"OK","level":"info","msg":"latency: 26.017µs","request":{"name":"user-2"},"response":{"message":"Hello user-2"},"start":"2021-10-17T16:34:56.497186+07:00","time":"2021-10-17T16:34:56+07:00","uri":"/hello.HelloService/Hello"}
{"STEP_1":"Hello service","client_ip":"127.0.0.1:62328","code":"OK","level":"info","msg":"latency: 23.865µs","request":{"name":"user-3"},"response":{"message":"Hello user-3"},"start":"2021-10-17T16:34:56.497793+07:00","time":"2021-10-17T16:34:56+07:00","uri":"/hello.HelloService/Hello"}
{"STEP_1":"Hello service","client_ip":"127.0.0.1:62328","code":"OK","level":"info","msg":"latency: 16.271µs","request":{"name":"user-4"},"response":{"message":"Hello user-4"},"start":"2021-10-17T16:34:56.498207+07:00","time":"2021-10-17T16:34:56+07:00","uri":"/hello.HelloService/Hello"}
{"STEP_1":"Hello service","client_ip":"127.0.0.1:62328","code":"OK","level":"info","msg":"latency: 17.859µs","request":{"name":"user-5"},"response":{"message":"Hello user-5"},"start":"2021-10-17T16:34:56.498731+07:00","time":"2021-10-17T16:34:56+07:00","uri":"/hello.HelloService/Hello"}
```
