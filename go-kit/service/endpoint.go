package main

import (
	"context"

	"github.com/ttkgof/microservices-examples/proto/greeter"

	"github.com/go-kit/kit/endpoint"
)

func newGreeterHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*greeter.HelloRequest)
		resp := &greeter.HelloResponse{}
		resp.Reply = "hello, " + req.Name
		return resp, nil
	}
}

func decodeHelloRequest(ctx context.Context, req interface{}) (interface{}, error) {
	return req, nil
}
func encodeHelloResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	return resp, nil
}

func newGreeterGoodbyeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*greeter.GoodbyeRequest)
		resp := &greeter.GoodbyeResponse{}
		resp.Reply = "goodbye, " + req.Name
		return resp, nil
	}
}

func decodeGoodbyeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	return req, nil
}
func encodeGoodbyeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	return resp, nil
}
