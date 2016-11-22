package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	cli "github.com/urfave/cli"
	"golang.org/x/net/context"
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
	app.Commands = []cli.Command{
		{
			Name:  "inspect",
			Usage: "Inspect the Health Check status of a container",
			Action: func(c *cli.Context) error {
				docker_client, client_err := createClient()
				if client_err != nil {
					panic(client_err)
				}
				log.Info("Connected to Docker daemon...")
				containerName := c.Args().First()
				log.Infof("Getting health for %s", containerName)
				containerJson, err := docker_client.ContainerInspect(context.Background(), containerName)
				if err != nil {
					if client.IsErrContainerNotFound(err) {
						fmt.Printf("Container '%s' not found", containerName)
						return nil
					} else {
						panic(err)
					}
				}
				log.Info(containerJson.State.Health.Status)
				return nil
			},
		},
		{
			Name:  "wait",
			Usage: "Wait until a container enters Healthy status",
			Action: func(c *cli.Context) error {
				panic("Not Implemented!")
			},
		},
	}
	app.Run(os.Args)
}

func createClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": user_agent}
	return client.NewClient("unix:///var/run/docker.sock", docker_api_version, nil, defaultHeaders)
}
