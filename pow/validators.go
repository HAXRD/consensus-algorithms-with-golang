package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
	"strconv"
)

type Validators struct {
	list []PublicKey // use each node/wallet's publicKey as identifier
}

func NewValidators(n int) *Validators {
	list := make([]PublicKey, n)
	for i := range n {
		list[i] = NewWallet("NODE-" + strconv.Itoa(i)).pubKey
	}
	return &Validators{list}
}

func (vs *Validators) ValidatorExists(validator PublicKey) bool {
	for _, pubKey := range vs.list {
		if pow_util.Byte2Hex(pubKey) == pow_util.Byte2Hex(validator) {
			return true
		}
	}
	return false
}
