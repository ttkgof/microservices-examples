package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ttkgof/microservices-examples/proto/greeter"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:19001", grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(requestIDInterceptor))
	if err != nil {
		fmt.Printf("connection error, %s\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	cli := greeter.NewGreeterClient(conn)

	resp, err := cli.Hello(context.Background(), &greeter.HelloRequest{Name: "xiaomage"})
	if err != nil {
		fmt.Printf("connection error, %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(resp.Reply)
}
