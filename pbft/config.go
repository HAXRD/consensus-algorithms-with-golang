package pbft

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TX_THRESHOLD int `yaml:"tx_threshold"` // Maximum num of txs for a block or can exist in the block pool
	NUM_OF_NODES int `yaml:"num_of_nodes"` // Num of nodes in the system
	MIN_APPROVAL int // Minimum num of positive votes for the message/block to be valid
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	config.MIN_APPROVAL = 2*(config.NUM_OF_NODES/3) + 1

	return &config, nil
}
