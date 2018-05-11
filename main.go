package main

import (
	"github.com/lsianturi/storest/app"
	"github.com/lsianturi/storest/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
