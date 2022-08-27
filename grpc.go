package logger

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

var DefaultConfigGrpc = ConfigGrpc{
	SkipperGrpc: DefaultSkipperGrpc,
}

func GrpcInterceptor(config ConfigGrpc) grpc.UnaryServerInterceptor {
	if config.SkipperGrpc == nil {
		config.SkipperGrpc = DefaultSkipperGrpc
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := New(WithFormatter(&JSONFormatter{}))
		if config.SkipperGrpc(ctx, info) {
			return handler(ctx, req)
		}

		if config.BeforeFuncGrpc != nil {
			config.BeforeFuncGrpc(ctx, info)
		}
		var (
			start   = time.Now()
			isError bool
		)

		p, ok := peer.FromContext(ctx)
		if ok {
			log.WithField(ClientIPField, p.Addr.String())
		}
		log.WithFields(map[string]interface{}{
			StartField:   start,
			URIField:     info.FullMethod,
			RequestField: req,
		})

		defer func() {
			var (
				code = codes.OK
				end  = time.Now()
			)

			if err != nil {
				code = status.Code(err)
			}
			log.WithField(CodeField, code.String())
			msg := fmt.Sprintf("latency: %v", end.Sub(start))
			if isError {
				log.Error(msg)
			} else {
				log.Info(msg)
			}
		}()

		ctx = context.WithValue(ctx, Key, log)
		resp, err = handler(ctx, req)
		if err != nil {
			isError = true
			log.WithField(ErrorsField, err)
		} else {
			log.WithField(ResponseField, resp)
		}
		return
	}
}
