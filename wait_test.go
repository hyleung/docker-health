package main

import (
	"github.com/docker/engine-api/types"
	"github.com/urfave/cli"
	"testing"
)

func TestWait_waitOnContainerHealth_no_healthcheck(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				Health: nil,
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	err := waitOnContainerHealth(stub, "foo", 1)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
}

func TestWait_waitOnContainerHealth_healthy(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "healthy",
				},
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	err := waitOnContainerHealth(stub, "foo", 1)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
}

func TestWait_waitOnContainerHealth_unhealthy(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "unhealthy",
				},
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	err := waitOnContainerHealth(stub, "foo", 1)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
}

func TestWait_waitOnContainerHealth_timeout(t *testing.T) {
	result := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				Health: &types.Health{
					Status: "starting",
				},
			},
		},
	}
	stub := StubContainerAPIClient{[]types.Container{}, nil, result, nil}
	err := waitOnContainerHealth(stub, "foo", 1)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	exitErr, ok := err.(cli.ExitCoder)
	if !ok {
		t.Fatal("Unexpected error type")
	}
	if exitErr.ExitCode() != 124 {
		t.Fatal("Unexpected error code: returned", exitErr.ExitCode(), ",but expected 124")
	}
}
