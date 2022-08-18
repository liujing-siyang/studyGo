package main

import (
	"fmt"
	"test_viper/config"
)

func main() {
	err := config.Initconfig()
	if err != nil {
		fmt.Printf("err %v", err)
	}
}
