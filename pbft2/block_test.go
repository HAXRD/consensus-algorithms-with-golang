package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestGenesis(t *testing.T) {
	genesis := Genesis()
	dHex := chain_util2.BytesToHex([]byte("------"))
	if chain_util2.BytesToHex(genesis.LastHash) != dHex ||
		chain_util2.BytesToHex(genesis.Hash) != dHex ||
		genesis.Data != nil ||
		genesis.Proposer != nil ||
		chain_util2.BytesToHex(genesis.Signature) != dHex ||
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
	hash := chain_util2.Hash(
		timestamp + chain_util2.BytesToHex(lastHash) + string(dataInByte) + strconv.FormatUint(nonce, 10),
	)
	if chain_util2.BytesToHex(hash) != chain_util2.BytesToHex(HashBlock(timestamp, lastHash, data, nonce)) {
		t.Error("HashBlock fail")
	}
}
