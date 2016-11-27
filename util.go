package main

import (
	"encoding/json"
	"github.com/docker/engine-api/types"
)

func toJson(data interface{}) string {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(result)
}

func Filter(list []types.ContainerJSON, f func(types.ContainerJSON) bool) []types.ContainerJSON {
	result := make([]types.ContainerJSON, 0)
	for _, container := range list {
		if f(container) {
			result = append(result, container)
		}
	}
	return result
}
