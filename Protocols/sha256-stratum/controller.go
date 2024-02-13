package sha256stratum

import (
	"fmt"
	"net"
	"strings"
)

// const net = require('net')
// const Miner = require('./miner')
// const logger = require('../../lib/logger')
// const config = require("../../config");

/*
TO DO

	Defines the Controller class which manages multiple miners and connections to various pools.
	Creates separate stratum proxy servers for each configured pool.
	Keeps track of connected miners and handles events like disconnections.
*/

// module.exports = class Controller {
//     constructor(coin, host, params) {
//         this.coin = coin
//         this.host = host
//         this.params = params
//         this.logPrefix = coin.toUpperCase() + ' >'

//         this.minersCount = 0
//         this.ports = {}
//         this.miners = new Map()

// }
type Controller struct {
	Coin        string
	Host        string
	Params      map[string]interface{}
	LogPrefix   string
	MinersCount int
	Ports       map[string]interface{}
	Miners      map[string]map[int]*Miner
}

func NewController(coin string, host string, params map[string]interface{}) Controller {
	return Controller{
		Coin:        coin,
		Host:        host,
		Params:      params,
		LogPrefix:   strings.ToUpper(coin) + " >",
		MinersCount: 0,
		Ports:       make(map[string]interface{}),
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
	c.Ports[port] = struct {
		Pool   interface{}
		Miners map[int]*Miner
		Server net.Listener
	}{
		Pool:   pool,
		Miners: make(map[int]*Miner),
		Server: server,
	}

	// Start accepting connections
	// go func() {
	// 	for {
	// 		// Accept a new connection
	// 		conn, err := server.Accept()
	// 		if err != nil {
	// 			// Handle error
	// 			continue
	// 		}

	// 		// Handle the connection in a new goroutine
	// 		go func(conn net.Conn) {
	// 			// Create a new Miner instance
	// 			miner := NewMiner(c, conn, port, pool)

	// 			// Connect to the pool
	// 			miner.ConnectToPool()

	// 			// Set fee (assuming fee is part of Controller's fields)
	// 			feeWindow := c.Params["fee"].(map[string]interface{})["window"].(int)
	// 			feePercent := c.Params["fee"].(map[string]interface{})["percent"].(float64)
	// 			miner.SetFee(feeWindow, feePercent)

	// 			// Add the miner to the map of miners for this port
	// 			c.Ports[port].Miners[miner.ID] = miner
	// 		}(conn)
	// 	}
	// }()

	// // Log that the Stratum proxy has started
	// fmt.Printf("Stratum proxy started (port: %s, pool: %v)\n", port, pool)
}

//     init() {
//         Object.keys(this.params.pools).forEach(port => {
//             this.createStratumProxy(port, this.params.pools[port])
//         })
//     }

//     createStratumProxy(port, pool) {
//         this.ports[port] = {
//             pool,
//             miners: new Map()
//         }
//         this.ports[port].server = net.createServer(socket => {
//             const miner = new Miner(this, socket, port, pool)
//             miner.connectToPool()
//             miner.setFee(this.params.fee.window, this.params.fee.percent)
//             this.ports[port].miners.set(miner.id, miner)
//         })
//         this.ports[port].server.listen(port, this.host, () => {
//             logger.info(this.logPrefix, `Stratum proxy started (port: ${port}, pool: ${pool})`)
//         })
//     }

//     /*
//     Close connection with ID connection(miner)
//      */
//     deleteMiner(port, id) {
//         this.ports[port].miners.delete(id)
//     }

//     /*
//     Close listening port
//      */
//     closePort(port) {
//         logger.info(this.logPrefix, `Сlose (port: ${port}, pool: ${this.ports[port].pool})`)
//         this.ports[port].server.close()
//         this.ports[port].miners.forEach(miner => miner.socket.destroy())
//         delete this.ports[port]
//     }

//     close() {
//         Object.keys(this.ports).forEach(port => {
//             logger.info(this.logPrefix, `Сlose (port: ${port}, pool: ${this.ports[port].pool})`)
//             this.ports[port].server.close()
//             this.ports[port].miners.forEach(miner => miner.socket.destroy())
//         })
//     }
// }
