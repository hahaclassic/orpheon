package main

import (
	cfgloader "github.com/hahaclassic/orpheon/pkg/config"
	"github.com/hahaclassic/orpheon/services/gateway-msv/config"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/app"
)

func main() {
	cfg := &config.Config{}
	cfgloader.Load(cfg)

	app.Run(cfg)
}
