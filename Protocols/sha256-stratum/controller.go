package sha256stratum

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
	Ports       map[string]string
	Miners      map[string]map[int]*Miner
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
