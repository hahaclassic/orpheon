package main

import (
	"github.com/hahaclassic/orpheon/services/user-msv/config"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/app"
)

func main() {
	cfg := config.Load()

	app.Run(cfg)
}
