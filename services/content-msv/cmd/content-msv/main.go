package main

import (
	"github.com/hahaclassic/orpheon/pkg/config"
	contentcfg "github.com/hahaclassic/orpheon/services/content-msv/config"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/app"
)

func main() {
	cfg := &contentcfg.Config{}
	config.Load(cfg)

	app.Run(cfg)
}
