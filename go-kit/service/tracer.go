package main

import (
	"context"
	"fmt"
	"io"

	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func initTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	config := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},

		ServiceName: serviceName,
	}

	var (
		tracer opentracing.Tracer
		closer io.Closer
		err    error
	)
	tracer, closer, err = config.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))

	if err != nil {
		return tracer, closer, err
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

func newTracer(tracer opentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName)
			defer span.Finish()
			return next(ctx, request)
		}
	}
}

type tracingKeyMD metadata.MD

func (md tracingKeyMD) ForeachKey(handler func(key, val string) error) error { //不能是指针
	for key, val := range md {
		for _, v := range val {
			if err := handler(key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func tracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("metadata create error")
		return handler(ctx, req)
	}
	md = md.Copy()
	tracer := opentracing.GlobalTracer()

	spanCtx, err := tracer.Extract(opentracing.TextMap, tracingKeyMD(md))
	if err != nil {
		fmt.Println("tracer extract error, ", err)
	}
	servSpan := tracer.StartSpan(info.FullMethod, opentracing.ChildOf(spanCtx))
	defer servSpan.Finish()

	ctx = opentracing.ContextWithSpan(ctx, servSpan)
	return handler(ctx, req)
}
