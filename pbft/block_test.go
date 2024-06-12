package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestGenesis(t *testing.T) {
	genesis := Genesis()
	dHex := chain_util.BytesToHex([]byte("------"))
	if chain_util.BytesToHex(genesis.LastHash) != dHex ||
		chain_util.BytesToHex(genesis.Hash) != dHex ||
		genesis.Data != nil ||
		chain_util.BytesToHex(genesis.Proposer) != dHex ||
		chain_util.BytesToHex(genesis.Signature) != dHex ||
		genesis.Nonce != 0 ||
		genesis.BlockMsgs != nil ||
		genesis.PrepareMsgs != nil ||
		genesis.CommitMsgs != nil ||
		genesis.RCMsgs != nil {
		t.Error("genesis block fail")
	}
}

func TestHashBlock(t *testing.T) {
	timestamp := time.Now().String()
	lastHash := []byte("-")
	var data []Transaction
	data = nil
	dataInByte, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	nonce := uint64(0)
	hash := chain_util.Hash(
		timestamp + chain_util.BytesToHex(lastHash) + string(dataInByte) + strconv.FormatUint(nonce, 10),
	)
	if chain_util.BytesToHex(hash) != chain_util.BytesToHex(HashBlock(timestamp, lastHash, data, nonce)) {
		t.Error("HashBlock fail")
	}
}
