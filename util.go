package main

import (
	"encoding/json"
)

func toJson(data interface{}) string {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(result)
}
