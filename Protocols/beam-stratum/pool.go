package beamstratum

import (
	"net"
)

type Pool struct {
	Miner        *Miner
	MiningServer string
	LogPrefix    string
	Socket       net.Conn
	NoncePrefix  *string
	Job          map[string]interface{}
	DataBuffer   string
}

func NewPool(miner *Miner) *Pool {
	return &Pool{
		Miner:        miner,
		MiningServer: miner.MiningServer,
		LogPrefix:    miner.LogPrefix,
		NoncePrefix:  nil,
		Job:          nil,
	}
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
