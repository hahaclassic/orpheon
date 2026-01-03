package main

import (
	"github.com/hahaclassic/orpheon/pkg/config"
	usercfg "github.com/hahaclassic/orpheon/services/user-msv/config"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/app"
)

func main() {
	cfg := &usercfg.Config{}
	config.Load(cfg)

	app.Run(cfg)
}
