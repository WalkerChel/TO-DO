package main

import "github.com/WalkerChel/TO-DO/internal/app"

const ConfigPath string = "configs/config.yaml"

func main() {
	app.Run(ConfigPath)
}
