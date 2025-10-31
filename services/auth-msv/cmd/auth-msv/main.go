package main

import (
	"github.com/hahaclassic/orpheon/pkg/config"
	authcfg "github.com/hahaclassic/orpheon/services/auth-msv/config"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/app"
)

func main() {
	cfg := &authcfg.Config{}
	config.Load(cfg)

	app.Run(cfg)
}
