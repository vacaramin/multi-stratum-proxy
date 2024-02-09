package main

import (
	"fmt"
	"multi-stratum-proxy/initializers"
)

func main() {
	config, err := initializers.ImportConfig("config.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Host:", config.Host)
	fmt.Println("Event Host:", config.Event.Host)
	fmt.Println("Event Port:", config.Event.Port)

	fmt.Println("Protocols:")
	for protocolName, protocolMap := range config.Protocols {
		fmt.Println("\tProtocol:", protocolName)
		for poolID, protocol := range protocolMap {
			fmt.Println("\t\tPool ID:", poolID)
			fmt.Println("\t\tPools:", protocol.Pools)
			fmt.Println("\t\tFee:")
			fmt.Println("\t\t\tWindow:", protocol.Fee.Window)
			fmt.Println("\t\t\tPercent:", protocol.Fee.Percent)
			fmt.Println("\t\t\tExtranonce_subscribe:", protocol.Fee.Extranonce_subscribe)
			fmt.Println("\t\t\tPool:", protocol.Fee.Pool)
			fmt.Println("\t\t\tWorker:", protocol.Fee.Worker)
			fmt.Println("\t\t\tPass:", protocol.Fee.Pass)
		}
	}
}
