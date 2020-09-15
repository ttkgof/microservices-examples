package main

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
)

//服务熔断
func newBreaker(cmd string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var resp interface{}
			err := hystrix.Do(cmd, func() (err error) {
				resp, err = next(ctx, request)
				return err
			}, nil)

			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}
