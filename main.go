package main

import (
	"fmt"
	"log"
	"multi-stratum-proxy/initializers"
)

func main() {
	config, err := initializers.Init("config.json")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(config.Host)
}
