package pbft

import "consensus-algorithms-with-golang/pbft/chainutil"

/**
Prepare is a type of message that stores passed-in blockHash, publicKey and signature.
It is used to guarantee FAULT TOLERANCE.
1. NewPrepare: create a prepare message with given blockHash, publicKey and signature
*/

type Prepare struct {
	blockHash string `json:"blockHash"`
	publicKey string `json:"publicKey"`
	signature string `json:"signature"`
}

func NewPrepare(wallet Wallet, block Block) *Prepare {
	return &Prepare{blockHash: block.hash, publicKey: chainutil.Key2Str(wallet.publicKey), signature: wallet.Sign(block.hash)}
}

/**
PreparePool stores a list of prepare message for "each" block

1. NewPreparePool: create a new prepare pool
2. InitPrepareForGivenBlock: init an empty list for a given blockHash
3. AddPrepare: pushes a prepare message for a blockhash into the list
4. PrepareExists: check if a given prepare message already exists
5. IsPrepareValid: check if the prepare message is valid or not
*/

type PreparePool struct {
	mapOfList map[string][]Prepare
}

func NewPreparePool() *PreparePool {
	return &PreparePool{make(map[string][]Prepare)}
}

func (pp *PreparePool) InitPrepareForGivenBlock(blockHash string) {
	pp.mapOfList[blockHash] = []Prepare{}
}

func (pp *PreparePool) AddPrepare2Pool(prepare Prepare) {
	pp.mapOfList[prepare.blockHash] = append(pp.mapOfList[prepare.blockHash], prepare)
}

func (pp *PreparePool) PrepareExists(prepare Prepare) bool {
	for _, p := range pp.mapOfList[prepare.blockHash] {
		if p.publicKey == prepare.publicKey {
			return true
		}
	}
	return false
}

func (pp *PreparePool) IsPrepareValid(prepare Prepare) bool {
	return chainutil.Verify(prepare.publicKey, prepare.blockHash, prepare.signature)
}
