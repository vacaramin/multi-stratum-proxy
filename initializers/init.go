package initializers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	beamstratum "multi-stratum-proxy/Protocols/beam-stratum"
	"net/http"
	"os"
	"reflect"
)

// Controller represents the controller for a specific mining protocol.
type Controller interface {
	Init(coin string, host string, params interface{})
	Close()
	ClosePort(port string)
	CreateStratumProxy(port string, pool interface{})
}

// Controllers maps protocol names to their respective controllers.
var Controllers = map[string]Controller{

	"beam-stratum": &beamstratum.Controller{},
	// "eth-proxy":        &EthProxyController{},
	// "ethereum-stratum": &EthereumStratumController{},
	// "grin-stratum":     &GrinStratumController{},
	// "kawpow-stratum":   &KawpowStratumController{},
	// "sha256-stratum":   &Sha256StratumController{},
	// "zhash-stratum":    &ZhashStratumController{},
	// "scrypt-stratum":    &ScryptStratumController{},
}

// Config represents the configuration for the proxy server.
var Config struct {
	Protocols map[string]map[string]interface{} `json:"protocols"`
	Event     struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"event"`
}

// InitProxy initializes the proxy server for a specific protocol and coin.
func InitProxy(protocol string, coin string, params interface{}) {
	controller, ok := Controllers[protocol]
	if !ok {
		fmt.Println("Invalid protocol specified:", protocol)
		return
	}
	controller.Init(coin, Config.Event.Host, params)
}

// CreateEventServer creates an event server for updating proxy configuration.
func CreateEventServer() {
	http.HandleFunc("/update-proxy-config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
			return
		}
		fmt.Println("EVENT > Update proxy config")
		ProxyUpdateEvent()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	addr := fmt.Sprintf("%s:%d", Config.Event.Host, Config.Event.Port)
	fmt.Printf("Started event server on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

// ProxyUpdateEvent updates proxy configuration based on received JSON data.
func ProxyUpdateEvent() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var updatedConfig struct {
		Protocols map[string]map[string]interface{} `json:"protocols"`
	}
	if err := json.Unmarshal(data, &updatedConfig); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if !reflect.DeepEqual(Config.Protocols, updatedConfig.Protocols) {
		// Perform necessary updates to proxy configurations.
	}

	Config.Protocols = updatedConfig.Protocols
}

func main() {
	// Read configuration from file.
	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	// Parse configuration.
	if err := json.Unmarshal(configData, &Config); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Initialize proxy for each protocol and coin.
	for protocol, coins := range Config.Protocols {
		for coin, params := range coins {
			InitProxy(protocol, coin, params)
		}
	}

	// Create event server for updating proxy configuration.
	CreateEventServer()
}
