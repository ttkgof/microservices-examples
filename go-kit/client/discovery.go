package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	register "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"

	"os"
	"sync"
	"time"
)

var ins *register.Instancer
var once sync.Once
var logger log.Logger
var pointer *sd.DefaultEndpointer

func newEndPointer() (*sd.DefaultEndpointer, error) {

	consulConfig := consul.DefaultConfig()
	consulConfig.Address = "0.0.0.0:8500"

	var apiCli *consul.Client
	apiCli, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	cli := register.NewClient(apiCli)
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "time", log.TimestampFormat(time.Now, "2006-01-02 15:04:05"))

	ins = register.NewInstancer(cli, logger, "ks", nil, true)

	pointer = sd.NewEndpointer(ins, greeterFactory, logger)

	return pointer, nil
}
