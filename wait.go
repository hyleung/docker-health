package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"time"
)

type WaitCommand struct {
}

func (*WaitCommand) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Wait on Healthcheck status for all running containers",
		},
		cli.Int64Flag{
			Name:  "timeout",
			Usage: "Wait timeout, in seconds",
			Value: 60,
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
		return waitOnContainerHealth(docker_client, containerName, c.Int64("timeout"))
	}
}

func waitOnContainerHealth(docker_client *client.Client, containerName string, timeout int64) error {
	log.Info(fmt.Sprintf("Waiting on health status of %s", containerName))
	_, err := docker_client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrContainerNotFound(err) {
			return cli.NewExitError(fmt.Sprintf("Container %s not found", containerName), 1)
		} else {
			return err
		}
	}
	timeout_channel := time.After(time.Duration(timeout) * time.Second)
	c := make(chan error, 1)
	go func() {
		for {
			containerJson, err := docker_client.ContainerInspect(context.Background(), containerName)
			if err != nil {
				if client.IsErrContainerNotFound(err) {
					c <- cli.NewExitError(fmt.Sprintf("Container %s not found", containerName), 1)
				} else {
					c <- err
				}
			}
			if containerJson.State.Health == nil {
				//If the container doesn't have health checks, exit normally
				log.Info(fmt.Sprintf("Container %s doesn't have any health checks defined", containerName))
				c <- nil
			} else if containerJson.State.Health.Status == "healthy" {
				log.Info(fmt.Sprintf("Container %s is healthy", containerName))
				//If the container is healthy, exit normally
				c <- nil
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	select {
	case <-timeout_channel:
		return cli.NewExitError(fmt.Sprintf("Container %s failed to enter healthy state after %d seconds", containerName, timeout), 1)
	case result := <-c:
		return result
	}
}
