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

func requestIDInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	rid := uuid.New().String()

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		d := md.Get(requestIDKey)
		if len(d) > 0 {
			rid = d[0]
		}
	}

	//这里必须放在OutgoingContext中才能传递给后面的
	ctx = metadata.AppendToOutgoingContext(ctx, requestIDKey, rid)
	ctx = context.WithValue(ctx, requestID{}, rid)

	fmt.Printf("set request id %s\n", rid)

	return invoker(ctx, method, req, reply, cc, opts...)
}
