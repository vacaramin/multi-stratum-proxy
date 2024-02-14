package sha256stratum

import (
	"log"
	"net"
	"strconv"
	"time"
)

/*
	TO DO

Defines the Miner class which represents a connected miner.
Establishes a connection to the pool proxy and handles incoming/outgoing messages.
Implements logic for handling fees by forwarding jobs to the fee pool and applying the calculated fee.
*/

type Fee struct {
	Pool     PoolToFee
	Active   bool
	Interval *time.Timer
	Timeout  *time.Timer
	Jobs     []map[string]string
}

type Miner struct {
	ID           int
	Controller   *Controller
	Socket       net.Conn
	Port         string
	MiningServer string
	LogPrefix    string
	Jobs         []map[string]interface{}
	Fee          Fee
}

func NewMiner(controller Controller, conn net.Conn, port string, miningServer string) *Miner {
	controller.MinersCount++
	log_prefix := controller.LogPrefix + " MINER:" + strconv.Itoa(controller.MinersCount-1) + " >"
	log.Printf("Miner Connecting (IP: %v ) ...", conn.RemoteAddr().String())

	return &Miner{
		ID:           controller.MinersCount - 1,
		Controller:   &controller,
		Socket:       conn,
		Port:         port,
		MiningServer: miningServer,
		LogPrefix:    log_prefix,
		Fee: Fee{
			Pool:     PoolToFee{},
			Active:   false,
			Interval: nil,
			Timeout:  nil,
			Jobs:     []map[string]string{},
		},
	}
}

func (m *Miner) ConnectToPool() {
	// Implementation for connecting to pool
}

func (m *Miner) SetFee(window int, percent int) {
	// Implementation for setting fee
}

func (m *Miner) SendJob(job []interface{}, fee bool, cleanJobs bool) {
	// Implementation for sending job
}

func (m *Miner) SendDifficulty(difficulty int, fee bool) {
	// Implementation for sending difficulty
}

func (m *Miner) SendVersionMask(mask string, fee bool) {
	// Implementation for sending version mask
}

func (m *Miner) SendExtranonce(extranonce1 string, extranonce2Size int, fee bool) {
	// Implementation for sending extranonce
}

func (m *Miner) Send(obj interface{}) {
	// Implementation for sending message to miner
}

func (m *Miner) SetEvents() {
	// Implementation for setting events for miner
}

func (m *Miner) HandleMessage(obj map[string]interface{}) {
	// Implementation for handling messages from miner
}
