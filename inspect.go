package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type InspectCommand struct {
}

type ContainerInfo struct {
	ID          string `json:"Id"`
	Image       string
	Name        string
	HealthCheck HealthCheck
}

type ContainerInfoShort struct {
	Image  string
	Name   string
	Status string
}

type HealthCheck struct {
	Command  string        `json:",omitempty"`
	Interval time.Duration `json:",omitempty"`
	Timeout  time.Duration `json:",omitempty"`
	Retries  int           `json:",omitempty"`
	Status   string
	Result   *types.HealthcheckResult
}

func (*InspectCommand) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Show Healthcheck status for all running containers",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show detailed health check information on containers",
		},
		cli.BoolFlag{
			Name:  "log, l",
			Usage: "Enable log output",
		},
	}

}
func (*InspectCommand) Command() interface{} {
	return func(c *cli.Context) error {
		if c.Bool("log") {
			log.SetLevel(log.InfoLevel)
		}
		docker_client, client_err := CreateClient()
		if client_err != nil {
			panic(client_err)
		}
		log.Info("Connected to Docker daemon...")
		if c.Bool("all") {
			result, err := healthForAllContainers(docker_client, c.Bool("verbose"))
			if err != nil {
				panic(err)
			}
			fmt.Println(result)
			return nil
		}
		containerName := c.Args().First()
		result, err := healthForContainer(docker_client, containerName, c.Bool("verbose"))
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		return nil
	}
}

func healthForContainer(docker_client DockerAPIClient, containerName string, verbose bool) (string, error) {
	log.Infof("Getting health for %s", containerName)
	containerJson, err := docker_client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrContainerNotFound(err) {
			return fmt.Sprintf("Container '%s' not found", containerName), nil
		} else {
			return "", err
		}
	}
	if containerJson.State.Health != nil {
		if verbose {
			result := ContainerInfo{
				ID:    containerJson.ID,
				Image: containerJson.Config.Image,
				Name:  containerJson.Name,
				HealthCheck: HealthCheck{Command: strings.Join(containerJson.Config.Healthcheck.Test, " "),
					Interval: containerJson.Config.Healthcheck.Interval,
					Timeout:  containerJson.Config.Healthcheck.Timeout,
					Retries:  containerJson.Config.Healthcheck.Retries,
					Status:   containerJson.State.Health.Status,
					Result:   containerJson.State.Health.Log[len(containerJson.State.Health.Log)-1],
				},
			}
			return toJson(result), nil
		} else {
			result := ContainerInfoShort{
				Name:   containerJson.Name,
				Image:  containerJson.Config.Image,
				Status: containerJson.State.Health.Status,
			}
			return toJson(result), nil
		}
	}
	return "{}", nil
}

func healthForAllContainers(docker_client DockerAPIClient, verbose bool) (string, error) {
	list, err := docker_client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", err
	}
	if verbose {
		containerJsonList := make([]ContainerInfo, 0)
		for _, v := range list {
			containerJson, _ := docker_client.ContainerInspect(context.Background(), v.ID)
			if containerJson.State.Health != nil {
				containerJsonList = append(containerJsonList, ContainerInfo{
					ID:    containerJson.ID,
					Image: containerJson.Config.Image,
					Name:  containerJson.Name,
					HealthCheck: HealthCheck{Command: strings.Join(containerJson.Config.Healthcheck.Test, " "),
						Interval: containerJson.Config.Healthcheck.Interval,
						Timeout:  containerJson.Config.Healthcheck.Timeout,
						Retries:  containerJson.Config.Healthcheck.Retries,
						Status:   containerJson.State.Health.Status,
						Result:   containerJson.State.Health.Log[len(containerJson.State.Health.Log)-1],
					},
				})
			}
		}
		return toJson(containerJsonList), nil
	} else {
		containerJsonList := make([]ContainerInfoShort, 0)
		for _, v := range list {
			containerJson, _ := docker_client.ContainerInspect(context.Background(), v.ID)
			if containerJson.State.Health != nil {
				containerJsonList = append(containerJsonList, ContainerInfoShort{
					Image:  containerJson.Config.Image,
					Name:   containerJson.Name,
					Status: containerJson.State.Health.Status,
				})
			}
		}
		return toJson(containerJsonList), nil
	}
}
