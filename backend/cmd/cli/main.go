package main

import (
	"github.com/hahaclassic/orpheon/backend/internal/cli"
	"github.com/hahaclassic/orpheon/backend/internal/config"
)

func main() {
	conf := config.MustLoad()

	cli.Run(conf)
}
