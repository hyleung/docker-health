package main

import (
	"errors"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
	"testing"
	"time"
)

type StubContainerAPIClient struct {
	Containers   []types.Container
	FetchListErr error
	Result       types.ContainerJSON
	Err          error
}

func (s StubContainerAPIClient) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	return s.Result, s.Err
}

func (s StubContainerAPIClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return s.Containers, s.FetchListErr
}

func TestInspect_healthForContainer_with_no_health(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				Health: nil,
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	r, err := healthForContainer(stub, "foo", false)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if r == nil {
		t.Error("Expected non-empty result")
	}
}

func TestInspect_healthForContainer_with_health(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Name: "some-name",
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "healthy",
				},
			},
		},
		Config: &container.Config{
			Image: "some-image",
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	r, err := healthForContainer(stub, "foo", false)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if r == nil {
		t.Error("Expected non-empty result")
	}
	_, ok := r.(ContainerInfo)
	if !ok {
		t.Error("Unexpected type")
	}
}

func TestInspect_healthForContainer_with_health_verbose(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Name: "some-name",
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "healthy",
					Log: []*types.HealthcheckResult{
						&types.HealthcheckResult{
							Start:    time.Time{},
							End:      time.Time{},
							ExitCode: 0,
							Output:   "some output",
						},
					},
				},
			},
		},
		Config: &container.Config{
			Image: "some-image",
			Healthcheck: &container.HealthConfig{
				Test:     []string{"/bin/sh"},
				Interval: 20,
				Timeout:  20,
				Retries:  1,
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	r, err := healthForContainer(stub, "foo", true)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if r == nil {
		t.Error("Expected non-empty result")
	}
	_, ok := r.(ContainerInfoDetail)
	if !ok {
		t.Error("Unexpected type")
	}
}

func TestInspect_healthForAllContainers(t *testing.T) {
	containerData := types.Container{
		ID: "some-container-id",
	}
	containerJson := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Name: "some-name",
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "healthy",
				},
			},
		},
		Config: &container.Config{
			Image: "some-image",
		},
	}
	stub := StubContainerAPIClient{[]types.Container{containerData}, nil, containerJson, nil}
	r, err := healthForAllContainers(stub, false)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if len(r) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestInspect_healthForAllContainers_verbose(t *testing.T) {
	containerData := types.Container{
		ID: "some-container-id",
	}
	containerJson := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Name: "some-name",
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "healthy",
					Log: []*types.HealthcheckResult{
						&types.HealthcheckResult{
							Start:    time.Time{},
							End:      time.Time{},
							ExitCode: 0,
							Output:   "some output",
						},
					},
				},
			},
		},
		Config: &container.Config{
			Image: "some-image",
			Healthcheck: &container.HealthConfig{
				Test:     []string{"/bin/sh"},
				Interval: 20,
				Timeout:  20,
				Retries:  1,
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{containerData}, nil, containerJson, nil}
	r, err := healthForAllContainers(stub, true)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if len(r) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestInspect_healthForAllContainers_inspect_error(t *testing.T) {
	containerData := types.Container{
		ID: "some-container-id",
	}
	stub := StubContainerAPIClient{[]types.Container{containerData}, nil, types.ContainerJSON{}, errors.New("ruh-roh")}
	_, err := healthForAllContainers(stub, false)
	if err == nil {
		t.Error("Expected error but non returned")
	}
}

func TestInspect_healthForAllContainers_inspect_error_vervose(t *testing.T) {
	containerData := types.Container{
		ID: "some-container-id",
	}
	stub := StubContainerAPIClient{[]types.Container{containerData}, nil, types.ContainerJSON{}, errors.New("ruh-roh")}
	_, err := healthForAllContainers(stub, true)
	if err == nil {
		t.Error("Expected error but non returned")
	}
}
