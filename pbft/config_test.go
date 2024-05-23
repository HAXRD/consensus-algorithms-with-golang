package pbft

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("config.yaml")

	if err != nil {
		t.Fatalf("Error loading config.yaml: %s", err)
	}
	if config.TX_THRESHOLD != 5 ||
		config.NUM_OF_NODES != 3 ||
		config.MIN_APPROVAL != 3 {
		t.Fatalf("LoadConfig did not load the correct configuration: %v, %v, %v",
			config.TX_THRESHOLD,
			config.NUM_OF_NODES,
			config.MIN_APPROVAL)
	}
}
