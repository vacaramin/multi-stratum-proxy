package beamstratum

import (
	"log"
	"net"
	"time"
)

// TO DO
// Defines two classes: Pool and PoolToFee.
// Pool connects to a mining pool via stratum protocol and handles communication with the pool.
// PoolToFee connects to a separate fee pool and forwards relevant job information for fee calculation.
type Pool struct {
	Miner        *Miner
	MiningServer string
	LogPrefix    string
	Socket       net.Conn
	NoncePrefix  *string
	Job          map[string]interface{}
	DataBuffer   string
}

// NewPool creates a new Pool instance(Similar to a constructor implementation)
func NewPool(miner *Miner) *Pool {
	return &Pool{
		Miner:        miner,
		MiningServer: miner.MiningServer,
		LogPrefix:    miner.LogPrefix,
		NoncePrefix:  nil,
		Job:          nil,
	}
}

// Connect is a method of Pool
func (p *Pool) Connect() {
	log.Printf("Connecting to pool: %s", p.MiningServer)

	socket, err := net.Dial("tcp", p.MiningServer)
	if err != nil {
		log.Printf("%s Error: %s", p.LogPrefix, err)
		return
	}
	p.Socket = socket
	socket.SetDeadline(time.Now().Add(10 * time.Second))

}

type PoolToFee struct {
	Miner       *Miner
	Params      map[string]interface{}
	LogPrefix   string
	Socket      net.Conn
	Status      string
	NoncePrefix string
	Job         map[string]interface{}
	DataBuffer  string
}
