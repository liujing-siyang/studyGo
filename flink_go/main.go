package main

import (
	"fmt"

	"github.com/flink-go/api"
)

func main() {
	// Your flink server HTTP API
	c, err := api.New("101.37.25.231:8081")
	if err != nil {
		panic(err)
	}

	// get cluster config
	config, err := c.Config()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}