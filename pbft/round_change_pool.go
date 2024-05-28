package pbft

import "consensus-algorithms-with-golang/pbft/chainutil"

/**
RoundChange is a type of message that stores passed-in blockHash, publicKey and signature,
along with a "ROUND-CHANGE" message.
1. NewRoundChange: create a round-change message with given blockHash, publicKey and signature.
*/

type RoundChange struct {
	blockHash string `json:"blockHash"`
	publicKey string `json:"publicKey"`
	signature string `json:"signature"`
	message   string `json:"message"`
}

func NewRoundChange(wallet Wallet, commit Commit) *RoundChange {
	return &RoundChange{
		commit.blockHash,
		chainutil.Key2Str(wallet.publicKey),
		wallet.Sign(commit.blockHash),
		"ROUND-CHANGE"}
}

/**
RoundChangePool stores a list of RoundChange message for "each" commit message

1. NetRoundChangePool: create a new RoundChange pool
2. InitRoundChangeForGivenCommit: init an empty list for a given blockHash
3. AddRoundChange2Pool: pushes a RoundChange message for a blockHash into the lis
4. RoundChangeExists: check if a given RoundChange message already exists
5. IsRoundChangeValid: check if the RoundChange message is valid or not
*/

type RoundChangePool struct {
	mapOfList map[string][]RoundChange
}

func NewRoundChangePool() *RoundChangePool {
	return &RoundChangePool{make(map[string][]RoundChange)}
}

func (rcp *RoundChangePool) InitRoundChangeForGivenCommit(blockHash string) {
	rcp.mapOfList[blockHash] = []RoundChange{}
}

func (rcp *RoundChangePool) AddRoundChange2Pool(roundChange RoundChange) {
	rcp.mapOfList[roundChange.blockHash] = append(rcp.mapOfList[roundChange.blockHash], roundChange)
}

func (rcp *RoundChangePool) RoundChangeExists(roundChange RoundChange) bool {
	for _, rc := range rcp.mapOfList[roundChange.blockHash] {
		if rc.blockHash == roundChange.blockHash {
			return true
		}
	}
	return false
}

func (rcp *RoundChangePool) IsRoundChangeValid(roundChange RoundChange) bool {
	return chainutil.Verify(roundChange.publicKey, roundChange.blockHash, roundChange.signature)
}
