package beamstratum

import (
	"net"
	"time"
)

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
