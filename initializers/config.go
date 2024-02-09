package initializers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Fee struct {
	Window               string
	Percent              int
	Extranonce_subscribe bool
	Pool                 string
	Worker               string
	Pass                 string
}

type Protocols struct {
	Protocol string
	Fee      Fee
}

type Event struct {
	Host string
	Port int
}

type Config struct {
	Host      string
	Protocols map[string]Protocols
	Event     Event
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
