package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	cli "github.com/urfave/cli"
	"os"
)

const (
	user_agent         = "engine-api-cli-1.0"
	docker_api_version = "1.24"
)

func main() {
	log.Debug("Starting docker-health...")
	app := cli.NewApp()
	app.Name = "docker-health"
	app.Usage = "Docker healthcheck utility"
	app.Version = "1.0"
	app.Action = func(c *cli.Context) error {
		defaultHeaders := map[string]string{"User-Agent": user_agent}
		_, client_err := client.NewClient("unix:///var/run/docker.sock", docker_api_version, nil, defaultHeaders)
		if client_err != nil {
			panic(client_err)
		}
		log.Info("Connected to Docker daemon...")
		return nil
	}
	app.Run(os.Args)
}
