package initializers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Fee struct {
	Window               int    `json:"window"`
	Percent              int    `json:"percent"`
	Extranonce_subscribe bool   `json:"extranonce_subscribe"`
	Pool                 string `json:"pool"`
	Worker               string `json:"worker"`
	Pass                 string `json:"pass"`
}

type Protocol struct {
	Pools map[string]string `json:"pools"`
	Fee   Fee               `json:"fee"`
}

type Protocols map[string]map[string]Protocol

type Event struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Config struct {
	Host      string                         `json:"host"`
	Protocols map[string]map[string]Protocol `json:"protocols"`
	Event     Event                          `json:"event"`
}

// ImportConfig : Unmarshal the config file using file name or path into the Config struct and return
func ImportConfig(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// PrintConfig : Function to test the unmarshalled Config struct
func PrintConfig(config Config) {
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
