package main

import (
	"github.com/hahaclassic/orpheon/backend/internal/app"
	"github.com/hahaclassic/orpheon/backend/internal/config"
)

func main() {
	conf := config.MustLoad()

	app.Run(conf)
}
