package main

import (
	"github.com/hahaclassic/orpheon/pkg/config"
	contentcfg "github.com/hahaclassic/orpheon/services/content-msv/config"
	app "github.com/hahaclassic/orpheon/services/content-msv/internal/app/http"
)

func main() {
	cfg := &contentcfg.Config{}
	config.Load(cfg)

	app.Run(cfg)
}
