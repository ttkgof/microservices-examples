package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/ttkgof/microservices-examples/proto/greeter"

	"google.golang.org/grpc"
)

//GreeterService greeter service
type GreeterService struct {
}

//Hello hello implementation
func (gs *GreeterService) Hello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloResponse, error) {
	resp := &greeter.HelloResponse{}
	fmt.Printf("hello from %s\n", req.Name)
	resp.Reply = "hello, " + req.Name
	return resp, nil
}

//Goodbye implementation
func (gs *GreeterService) Goodbye(ctx context.Context, req *greeter.GoodbyeRequest) (*greeter.GoodbyeResponse, error) {
	resp := &greeter.GoodbyeResponse{}
	fmt.Printf("goodbye from %s\n", req.Name)
	resp.Reply = "goodbye, " + req.Name
	return resp, nil
}

func main() {
	s := &GreeterService{}

	ls, err := net.Listen("tcp", "0.0.0.0:19001")
	if err != nil {
		fmt.Printf("net listen error, %s\n", err.Error())
		os.Exit(1)
	}

	gs := grpc.NewServer(grpc.ChainUnaryInterceptor(requestIDInterceptor))
	greeter.RegisterGreeterService(gs, greeter.NewGreeterService(s))
	gs.Serve(ls)
}
