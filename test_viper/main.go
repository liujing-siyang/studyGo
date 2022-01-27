package main

import (
	"fmt"
	"test_viper/config"
)

func main() {
	err := config.Initconfig1()
	if err != nil {
		fmt.Printf("err %v", err)
	}
}
