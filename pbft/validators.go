package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
	"strconv"
)

/**
PBFT is a permission required blockchain consensus algorithm.
Therefore, only nodes in the validator list are considered valid.
Our validators struct features the following methods:
1. NewValidators
2. ValidatorExists

NOTE:
The secret to create a key-pair is usually a 128/256-bit seed
generated from a source of entropy.
However, for demonstration purpose, we use `NODE-{i}` as the secret
for keypair generation. It's worth noting that in practice, the
secret should never be revealed publicly, otherwise the security of
the wallet is compromised.
*/

type Validators struct {
	list []PublicKey // use each node/wallet's publicKey as identifier
}

// NewValidators creates a slice of given num of validators
func NewValidators(n int) *Validators {
	list := make([]PublicKey, n)
	for i := range n {
		list[i] = NewWallet("NODE-" + strconv.Itoa(i)).publicKey
	}
	return &Validators{list}
}

// ValidatorExists checks if a node/wallet is within the list
func (vs *Validators) ValidatorExists(validator PublicKey) bool {
	for _, pubKey := range vs.list {
		if chain_util.BytesToHex(pubKey) == chain_util.BytesToHex(validator) {
			return true
		}
	}
	return false
}
