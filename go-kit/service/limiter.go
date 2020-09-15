package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

//服务限流， 令牌桶算法
func newDelayingLimiter(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			err := limit.Wait(ctx)
			if err != nil {
				return nil, err
			}
			return next(ctx, request)
		}
	}
}
