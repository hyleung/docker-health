package main

import (
	"github.com/docker/engine-api/types"
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
	err := waitOnContainerHealth(stub, "foo", 5)
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
	err := waitOnContainerHealth(stub, "foo", 5)
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
	err := waitOnContainerHealth(stub, "foo", 5)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
}
