package initializers

import "fmt"

func Init(filename string) (*Config, error) {
	fmt.Println("Initializers initialized")
	Config, err := ImportConfig(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return Config, nil
}
