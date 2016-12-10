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

func init() {
	log.SetLevel(log.WarnLevel)
}

func main() {
	log.Debug("Starting docker-health...")
	app := cli.NewApp()
	app.Name = "docker-health"
	app.Usage = "Docker healthcheck utility"
	app.Version = "1.0"

	inspectCommand := InspectCommand{}
	waitCommand := WaitCommand{}
	app.Commands = []cli.Command{
		{
			Name:   "inspect",
			Usage:  "Inspect the Health Check status of a container",
			Flags:  inspectCommand.Flags(),
			Action: inspectCommand.Command(),
		},
		{
			Name:   "wait",
			Usage:  "Wait until a container enters Healthy status",
			Flags:  waitCommand.Flags(),
			Action: waitCommand.Command(),
		},
	}
	app.Run(os.Args)
}

func CreateClient() (*client.Client, error) {
	if os.Getenv("DOCKER_HOST") != "" {
		return client.NewEnvClient()
	} else {
		defaultHeaders := map[string]string{"User-Agent": user_agent}
		return client.NewClient("unix:///var/run/docker.sock", docker_api_version, nil, defaultHeaders)
	}
}
