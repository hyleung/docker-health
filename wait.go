package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"strings"
)

type WaitCommand struct {
}

func (*WaitCommand) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Wait on Healthcheck status for all running containers",
		},
	}
}

func (*WaitCommand) Command() interface{} {
	return func(c *cli.Context) error {
		docker_client, client_err := CreateClient()
		if client_err != nil {
			panic(client_err)
		}
		log.Debug("Connected to Docker daemon...")
		if c.Bool("all") {
			log.Info("Not implemented....")
			return cli.NewExitError("Ruh-roh!", 1)
		}
		containerName := c.Args().First()
		return waitOnContainerHealth(docker_client, containerName)
	}
}

func waitOnContainerHealth(docker_client *client.Client, containerName string) error {
	log.Info(fmt.Sprintf("Waiting on health status of %s", containerName))
	containerJson, err := docker_client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrContainerNotFound(err) {
			return cli.NewExitError(fmt.Sprintf("Container %s not found", containerName), 1)
		} else {
			panic(err)
		}
	}
	if containerJson.State.Health == nil {
		//If the container doesn't have health checks, exit normally
		log.Info(fmt.Sprintf("Container %s doesn't have any health checks defined", containerName))
		return nil
	} else if containerJson.State.Health.Status == "healthy" {
		log.Info(fmt.Sprintf("Container %s is healthy", containerName))
		//If the container is healthy, exit normally
		return nil
	} else {
		//Keep checking
		args := filters.NewArgs()
		args.Add("Type", events.DaemonEventType)
		readCloser, err := docker_client.Events(context.Background(), types.EventsOptions{Filters: args})
		defer readCloser.Close()
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(readCloser)
		for scanner.Scan() {
			var event events.Message
			err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&event)
			action := event.Action
			actor := event.Actor
			if actor.Attributes["name"] == containerName {
				log.Info("Action: %s, Actor: %s", action, actor)
			}
			return nil
		}
		return nil
	}
}
