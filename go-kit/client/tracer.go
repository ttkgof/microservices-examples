package main

import (
	"context"
	"fmt"
	"io"
	"strings"

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

	tracer, closer, err := config.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
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

//参考 https://blog.csdn.net/liyunlong41/article/details/88043604
type tracingKeyMD metadata.MD

func (md tracingKeyMD) Set(k string, val string) {
	k = strings.ToLower(k)
	md[k] = append(md[k], val)
}

func tracingInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	var pCtx opentracing.SpanContext
	pSpan := opentracing.SpanFromContext(ctx)
	if pSpan != nil {
		pCtx = pSpan.Context()
	}
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan(method, opentracing.ChildOf(pCtx))
	defer span.Finish()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("metadata create error")
	}
	//参考FromIncomingContext注释，修改md前需要Copy
	md = md.Copy()
	err := tracer.Inject(span.Context(), opentracing.TextMap, tracingKeyMD(md))
	if err != nil {
		fmt.Println("tracer inject error, ", err)
	}

	//创建一个新的context，把metadata附带上
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return invoker(newCtx, method, req, reply, cc, opts...)
}
