package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	register "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
)

//newRegister 服务注册
func newRegister(port string) (sd.Registrar, error) {
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = "0.0.0.0:8500"

	c, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	cli := register.NewClient(c)

	check := consul.AgentServiceCheck{
		TCP: "192.168.50.141:" + port,
		//HTTP:     "http://" + "127.0.0.1:" + "19002" + "/health",
		Interval: "3s",
		//TTL: "9s",
		Timeout:                        "15s",
		Notes:                          "health check",
		DeregisterCriticalServiceAfter: "1m",
	}
	p, _ := strconv.Atoi(port)
	rand.Seed(time.Now().UTC().UnixNano())
	num := rand.Intn(100)
	asr := consul.AgentServiceRegistration{
		ID:      fmt.Sprintf("ks_%d", num),
		Name:    "ks",
		Address: "127.0.0.1",
		Port:    p,
		Tags:    nil,
		Check:   &check,
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "time", log.TimestampFormat(time.Now, "2006-01-02 15:04:05"))
	return register.NewRegistrar(cli, &asr, logger), nil
}
