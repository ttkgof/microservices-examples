package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
)

func newBreaker() endpoint.Middleware {
	cmd := "ks"
	hystrix.ConfigureCommand(cmd, hystrix.CommandConfig{
		Timeout:                1000 * 30, //超时时间 毫秒
		ErrorPercentThreshold:  1,         //错误率 请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动
		SleepWindow:            10000,     //休眠时间 毫秒
		MaxConcurrentRequests:  1000,      //最大并发量
		RequestVolumeThreshold: 5,         //请求阈值(一个统计窗口10秒内请求数量)  熔断器是否打开首先要满足这个条件；这里的设置表示至少有5个请求才进行ErrorPercentThreshold错误百分比计算
	})

	return circuitbreaker.Hystrix(cmd)
}
