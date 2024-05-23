package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"strconv"
)

/**
PBFT is a permission required blockchain consensus algorithm. Therefore,
only nodes in the validator list are considered as valid validator.

NOTE:
The secret to create keypairs are usually a 128/256-bit seeds generated
from a source of entropy.
However, for demonstration purpose, we use `NODE-{i}` as the secret for keypair
generation. It's worth noting that in practice, the secret should never be revealed
publicly, otherwise the security of the wallet is compromised.

1. generate a list of nodes as validators (use wallet's publicKey)
2. check if a given node's publicKey is valid (is one of hte validators)
*/

type Validators struct {
	list []string // use each node/wallet's publicKey as identifier
}

// NewValidators creates a slice of given num of validators
func NewValidators(numOfValidators int) *Validators {
	list := make([]string, numOfValidators)
	for i := 0; i < numOfValidators; i++ {
		list[i] = chainutil.Key2Str(NewWallet("NODE-" + strconv.Itoa(i)).GetPublicKey())
	}
	return &Validators{list}
}

// IsAValidator checks if a node/wallet is a validator
func (v *Validators) IsAValidator(validator string) bool {
	for _, publicKey := range v.list {
		if publicKey == validator {
			return true
		}
	}
	return false
}
