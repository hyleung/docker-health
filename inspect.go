package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func InspectContainerCommand() interface{} {
	return func(c *cli.Context) error {
		docker_client, client_err := CreateClient()
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
	}
}
