package main

import (
	"mc-tool/application"
	"mc-tool/cli"
)

func main() {
	cli.Run(application.NewApplication())
}
