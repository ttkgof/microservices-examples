package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ttkgof/microservices-examples/proto/greeter"

	"golang.org/x/time/rate"

	transport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"github.com/google/uuid"
)

type requestIDKey struct{}
type methodKey struct{}

//GreeterService greeter service
type GreeterService struct {
	helloHandler   transport.Handler
	goodbyeHandler transport.Handler
}

//Hello hello implementation
func (gs *GreeterService) Hello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloResponse, error) {
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05.9"), "hello", ctx)
	tmpCtx := context.WithValue(ctx, requestIDKey{}, uuid.New().String())
	tmpCtx = context.WithValue(tmpCtx, methodKey{}, "hello")
	_, resp, err := gs.helloHandler.ServeGRPC(tmpCtx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*greeter.HelloResponse), nil
}

//Goodbye implementation
func (gs *GreeterService) Goodbye(ctx context.Context, req *greeter.GoodbyeRequest) (*greeter.GoodbyeResponse, error) {

	//fmt.Println(time.Now().Format("2006-01-02 15:04:05.9"), "goodbye", ctx)
	tmpCtx := context.WithValue(ctx, requestIDKey{}, uuid.New().String())
	tmpCtx = context.WithValue(tmpCtx, methodKey{}, "goodbye")

	_, resp, err := gs.goodbyeHandler.ServeGRPC(tmpCtx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*greeter.GoodbyeResponse), nil
}

func main() {
	var port string
	flag.StringVar(&port, "p", "19001", "rpc port")
	flag.Parse()

	logger := initLogger()

	tracer, tracingCloser, err := initTracer("ks")
	defer tracingCloser.Close()
	if err != nil {
		fmt.Println("init tracer error ", err)
		os.Exit(1)
	}

	s := &GreeterService{}

	limitCount := 1000
	limiter := rate.NewLimiter(rate.Every(time.Second*1), limitCount)

	greeterHelloEndPoint := newGreeterHelloEndpoint()
	greeterHelloEndPoint = newDelayingLimiter(limiter)(greeterHelloEndPoint)
	greeterHelloEndPoint = newBreaker("greeter")(greeterHelloEndPoint)
	greeterHelloEndPoint = newLogger(logger)(greeterHelloEndPoint)
	greeterHelloEndPoint = newTracer(tracer, "hello")(greeterHelloEndPoint)

	greeterGoodbyeEndPoint := newGreeterGoodbyeEndpoint()
	greeterGoodbyeEndPoint = newDelayingLimiter(limiter)(greeterGoodbyeEndPoint)
	greeterGoodbyeEndPoint = newBreaker("greeter")(greeterGoodbyeEndPoint)
	greeterGoodbyeEndPoint = newLogger(logger)(greeterGoodbyeEndPoint)
	greeterGoodbyeEndPoint = newTracer(tracer, "goodbye")(greeterGoodbyeEndPoint)

	s.helloHandler = transport.NewServer(
		greeterHelloEndPoint,
		decodeHelloRequest,
		encodeHelloResponse,
	)
	s.goodbyeHandler = transport.NewServer(
		greeterGoodbyeEndPoint,
		decodeGoodbyeRequest,
		encodeGoodbyeResponse,
	)

	ls, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Printf("net listen error, %s\n", err.Error())
		os.Exit(1)
	}

	reg, err := newRegister(port)
	if err != nil {
		fmt.Printf("net listen error, %s\n", err.Error())
		os.Exit(1)
	}

	reg.Register()

	//defer reg.Deregister()

	//grpc并不能传递context,这里需要使用middleware改成http header，然后后面解析
	// interceptor := grpc_middleware.ChainUnaryServer(
	// 	grpc_opentracing.Unar
	// 	yServerInterceptor(
	// 		grpc_opentracing.WithTracer(tracer)),
	// )
	gs := grpc.NewServer(grpc.UnaryInterceptor(tracingInterceptor))
	greeter.RegisterGreeterService(gs, greeter.NewGreeterService(s))

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		c := <-quit
		fmt.Println("signal:", c)

		reg.Deregister()
		logger.Sync()
		tracingCloser.Close()

		gs.Stop()
	}()

	fmt.Println("start")
	err = gs.Serve(ls)
	if err != nil {
		fmt.Println(err)
	}
}
