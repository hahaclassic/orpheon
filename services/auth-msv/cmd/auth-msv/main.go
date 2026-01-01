package main

import (
	"github.com/hahaclassic/orpheon/pkg/config"
	authcfg "github.com/hahaclassic/orpheon/services/auth-msv/config"
	app "github.com/hahaclassic/orpheon/services/auth-msv/internal/app/http"
)

func main() {
	cfg := &authcfg.Config{}
	config.Load(cfg)

	app.Run(cfg)
}
