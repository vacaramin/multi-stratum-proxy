package main

import (
	"fmt"
	"multi-stratum-proxy/initializers"
)

/*
Implementation
Starts the proxy server based on the provided configuration file.
Parses the configuration and creates instances of Controller for different protocols (beam-stratum in this case).
Creates an event server to listen for updates to the configuration file.
Reloads the configuration and updates running proxies whenever the configuration changes.
*/
func main() {
	/*config*/ _, err := initializers.ImportConfig("config.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//initializers.PrintConfig(*config)
}
