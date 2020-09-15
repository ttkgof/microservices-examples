package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const requestIDKey = "X-Request-ID"

type requestID struct{}

func requestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	rid := uuid.New().String()

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		d := md.Get(requestIDKey)
		if len(d) > 0 {
			rid = d[0]
		}
	}
	ctx = metadata.AppendToOutgoingContext(ctx, requestIDKey, rid)
	//方便后面的使用
	ctx = context.WithValue(ctx, requestID{}, rid)
	fmt.Printf("set request id %s\n", rid)
	return handler(ctx, req)
}
