package beamstratum

import (
	"fmt"
	"net"
	"strconv"
)

type Controller struct {
	Coin      string
	Host      string
	Params    interface{}
	LogPrefix string

	MinersCount int
	Ports       map[string]*Port
	Miners      map[string]*Miner
}

type Port struct {
	Pool   interface{} // Update the type accordingly
	Miners map[string]*Miner
	Server net.Listener
}

func NewController(coin, host string, params map[string]interface{}) *Controller {
	return &Controller{
		Coin:        coin,
		Host:        host,
		Params:      params,
		LogPrefix:   coin + " >",
		MinersCount: 0,
		Ports:       make(map[string]*Port),
		Miners:      make(map[string]*Miner),
	}
}
func (c *Controller) Init(coin, host string, params interface{}) {
	c.Coin = coin
	c.Host = host
	c.Params = params
	c.LogPrefix = coin + " >"
	c.MinersCount = 0
	c.Ports = make(map[string]*Port)
	c.Miners = make(map[string]*Miner)
}

func (c *Controller) CreateStratumProxy(port string, pool interface{}) {
	c.Ports[port] = &Port{
		Pool:   pool,
		Miners: make(map[string]*Miner),
	}

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%s", c.Host, port))
	if err != nil {
		// Handle error
		return
	}

	c.Ports[port].Server = server

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				// Handle error
				continue
			}

			miner := NewMiner(c, conn, port, pool)
			miner.ConnectToPool()
			// SetFee logic here

			c.Ports[port].Miners[strconv.Itoa(miner.ID)] = miner
		}
	}()
}

func (c *Controller) DeleteMiner(port, id string) {
	delete(c.Ports[port].Miners, id)
}

func (c *Controller) ClosePort(port string) {
	fmt.Printf("%s Close (port: %s, pool: %v)\n", c.LogPrefix, port, c.Ports[port].Pool)
	c.Ports[port].Server.Close()

	for _, miner := range c.Ports[port].Miners {
		miner.Socket.Close()
	}
	delete(c.Ports, port)
}

func (c *Controller) Close() {
	for port := range c.Ports {
		c.ClosePort(port)
	}
}
