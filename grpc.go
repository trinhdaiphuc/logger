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

func GrpcInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var (
		log       = New(WithFormatter(&JSONFormatter{}))
		start     = time.Now()
		panicking = true
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

		switch {
		case err != nil:
			code = status.Code(err)
		case panicking:
			code = codes.Internal
		}
		log.WithField(CodeField, code.String())
		msg := fmt.Sprintf("latency: %v", end.Sub(start))
		log.Info(msg)
	}()

	ctx = context.WithValue(ctx, Key, log)
	resp, err = handler(ctx, req)
	panicking = false // normal exit, no panic happened, disarms defer
	log.WithField(ResponseField, resp)
	return
}
