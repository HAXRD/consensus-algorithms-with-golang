package pbft

import (
	"testing"
)

// TestNewValidators test if NewValidators generates given num of validators
func TestNewValidators(t *testing.T) {
	validators := *NewValidators(5)

	if 5 != len(validators.list) {
		t.Errorf("NewValidators failed, expected 5, got %d", len(validators.list))
	}
}

// TestIsAValidator
func TestIsAValidator(t *testing.T) {
	validators_a := *NewValidators(5)
	validators_b := *NewValidators(7)

	if !validators_a.IsAValidator(validators_a.list[0]) {
		t.Errorf("IsAValidator failed, expected true, got false")
	}
	if validators_a.IsAValidator(validators_b.list[6]) {
		t.Errorf("IsAValidator failed, expected false, got true")
	}
}
