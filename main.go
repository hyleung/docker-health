package main

import (
	cli "github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-health"
	app.Usage = "Docker healthcheck utility"
	app.Version = "1.0"
	app.Run(os.Args)
}
