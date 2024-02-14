package sha256stratum

import (
	"fmt"
	"net"
	"strings"
)

type PortDetail struct {
	Pool   interface{}
	Miners map[int]*Miner
	Server net.Listener
}

type Controller struct {
	Coin        string
	Host        string
	Params      map[string]interface{}
	LogPrefix   string
	MinersCount int
	Ports       map[string]PortDetail
	Miners      map[string]map[int]*Miner
}

func NewController(coin string, host string, params map[string]interface{}) Controller {
	return Controller{
		Coin:        coin,
		Host:        host,
		Params:      params,
		LogPrefix:   strings.ToUpper(coin) + " >",
		MinersCount: 0,
		Ports:       make(map[string]PortDetail),
		Miners:      make(map[string]map[int]*Miner),
	}
}
func (c *Controller) Init() {
	for port, pool := range c.Params["pools"].(map[string]interface{}) {
		c.CreateStratumProxy(port, pool)
	}
}

func (c *Controller) CreateStratumProxy(port string, pool interface{}) {
	// Create a TCP server
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%s", c.Host, port))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new map entry for the port
	c.Ports[port] = PortDetail{
		Pool:   pool,
		Miners: make(map[int]*Miner),
		Server: server,
	}

	fmt.Printf("Stratum proxy started (port: %s, pool: %v)\n", port, pool)
}

func (c *Controller) DeleteMiner(port string, id int) {
	delete(c.Miners[port], id)
}

func (c *Controller) ClosePort(port string) {
	if _, ok := c.Ports[port]; ok {
		c.Ports[port].Server.Close()
		delete(c.Ports, port)
	}
}

func (c *Controller) Close() {
	for port := range c.Ports {
		c.ClosePort(port)
	}
}
