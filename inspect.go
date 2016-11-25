package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

type InspectCommand struct {
}

func (*InspectCommand) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Show Healthcheck status for all running containers",
		},
	}

}
func (*InspectCommand) Command() interface{} {
	return func(c *cli.Context) error {
		docker_client, client_err := CreateClient()
		if client_err != nil {
			panic(client_err)
		}
		log.Debug("Connected to Docker daemon...")
		if c.Bool("all") {
			healthForAllContainers(docker_client)
			return nil
		}
		containerName := c.Args().First()
		healthForContainer(docker_client, containerName)
		return nil
	}
}

func healthForContainer(docker_client *client.Client, containerName string) {
	log.Debugf("Getting health for %s", containerName)
	containerJson, err := docker_client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrContainerNotFound(err) {
			fmt.Printf("Container '%s' not found", containerName)
		} else {
			panic(err)
		}
	}
	if containerJson.State.Health != nil {
		b, _ := json.MarshalIndent(containerJson.State.Health, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Println("{}")
	}
}

func healthForAllContainers(docker_client *client.Client) {
	fmt.Println("To be continued...")
}
