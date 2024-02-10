package beamstratum

/*
TO DO

	Defines the Controller class which manages multiple miners and connections to various pools.
	Creates separate stratum proxy servers for each configured pool.
	Keeps track of connected miners and handles events like disconnections.
*/
type Controller struct {
	Coin        string
	Host        string
	Params      map[string]interface{}
	LogPrefix   string
	MinersCount int
	Ports       map[string]string
	Miners      map[string]map[int]*Miner
}
