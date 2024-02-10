package beamstratum

type Controller struct {
	Coin        string
	Host        string
	Params      map[string]interface{}
	LogPrefix   string
	MinersCount int
	Ports       map[string]string
	Miners      map[string]map[int]*Miner
}
