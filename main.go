package main

import (
	"fmt"
	"multi-stratum-proxy/initializers"
)

func main() {
	/*config*/ _, err := initializers.ImportConfig("config_19.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//initializers.PrintConfig(*config)
}
