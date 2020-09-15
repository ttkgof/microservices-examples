package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/ttkgof/microservices-examples/proto/greeter"
	"google.golang.org/grpc"
)

func greeterFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		fmt.Println(instanceAddr)
		c, err := grpc.Dial(instanceAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(tracingInterceptor))
		if err != nil {
			return nil, err
		}
		defer c.Close()
		cli := greeter.NewGreeterClient(c)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.9"), ctx)
		switch r := req.(type) {
		case *greeter.HelloRequest:
			return cli.Hello(ctx, r)
		case *greeter.GoodbyeRequest:
			return cli.Goodbye(ctx, r)
		case greeter.HelloRequest:
			return cli.Hello(ctx, &r)
		case greeter.GoodbyeRequest:
			return cli.Goodbye(ctx, &r)
		default:
			return nil, errors.New("unknown request type")
		}
	}, nil, nil
}
