package main

import (
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

type DockerAPIClient interface {
	ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
}
