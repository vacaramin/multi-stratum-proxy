package beamstratum

import (
	"net"
	"time"
)

/*
	TO DO

Defines the Miner class which represents a connected miner.
Establishes a connection to the pool proxy and handles incoming/outgoing messages.
Implements logic for handling fees by forwarding jobs to the fee pool and applying the calculated fee.
*/

type Miner struct {
	ID           int
	Controller   *Controller
	Socket       net.Conn
	Port         string
	MiningServer string
	LogPrefix    string
	Jobs         []map[string]interface{}
	Fee          struct {
		Pool     PoolToFee
		Active   bool
		Interval *time.Timer
		Timeout  *time.Timer
		Jobs     []map[string]string
	}
}
