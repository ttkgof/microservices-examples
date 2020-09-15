package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/ttkgof/microservices-examples/proto/greeter"
)

func sayHello(pointer *sd.DefaultEndpointer) {
	//负载均衡
	balancer := lb.NewRoundRobin(pointer)
	requestClient := lb.Retry(3, 3*time.Second, balancer)

	tracer, closer, err := initTracer("kc")
	if err != nil {
		fmt.Printf(time.Now().Format("2006-01-02 15:04:05.9"), err)
	}

	defer closer.Close()

	requestClient = newBreaker()(requestClient)
	requestClient = newTracer(tracer, "hello")(requestClient)
	req := &greeter.HelloRequest{Name: "xiaomage"}

	resp, err := requestClient(context.Background(), req)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.9"), resp, err)
}
func sayGoodbye(pointer *sd.DefaultEndpointer) {
	balancer := lb.NewRoundRobin(pointer)
	requestClient := lb.Retry(3, 3*time.Second, balancer)

	tracer, closer, err := initTracer("kc")
	if err != nil {
		fmt.Printf(time.Now().Format("2006-01-02 15:04:05.9"), err)
	}

	defer closer.Close()

	requestClient = newBreaker()(requestClient)
	requestClient = newTracer(tracer, "goodbye")(requestClient)
	req := &greeter.GoodbyeRequest{Name: "xiaomage"}

	resp, err := requestClient(context.Background(), req)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.9"), resp, err)
}

func main() {
	//服务发现
	p, err := newEndPointer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, closer, err := initTracer("ks")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closer.Close()
	//测试限流
	//for i := 0; i < 30; i++ {
	//	sayHello(p)
	//	//测试熔断，中途断开服务器会显示hystrix: circuit open
	//	time.Sleep(time.Second * 1)
	//}
	//time.Sleep(time.Second * 30)
	sayGoodbye(p)
}
