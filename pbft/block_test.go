package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"encoding/json"
	"testing"
)

func TestNewBlock(t *testing.T) {
	timestamp := "timestamp"
	lastHash := "lastHash"
	hash := "hash"
	data := []Transaction{}
	proposer := "proposer"
	signature := "signature"
	sequenceNum := uint64(0)

	newBlock := NewBlock(timestamp, lastHash, hash, data, proposer, signature, sequenceNum)
	if timestamp != newBlock.timestamp ||
		lastHash != newBlock.lastHash ||
		hash != newBlock.hash ||
		len(data) != len(newBlock.data) ||
		proposer != newBlock.proposer ||
		signature != newBlock.signature ||
		sequenceNum != newBlock.sequenceNum {
		t.Errorf("NewBlock fail")
	}
	newTimestamp := "new_timestamp"
	newBlock.timestamp = newTimestamp
	if newTimestamp != newBlock.timestamp {
		t.Errorf("NewBlock fail")
	}
}

func TestGenesis(t *testing.T) {
	genesis := Genesis()
	if "genesis time" != genesis.timestamp ||
		"----" != genesis.lastHash ||
		"genesis-hash" != genesis.hash ||
		0 != len(genesis.data) ||
		"====" != genesis.proposer ||
		"SIGN" != genesis.signature ||
		uint64(0) != genesis.sequenceNum {
		t.Errorf("Genesis fail")
	}
	newTimestamp := "new_timestamp"
	genesis.timestamp = newTimestamp
	if newTimestamp != genesis.timestamp {
		t.Errorf("NewBlock fail")
	}
}

func TestHashBlock(t *testing.T) {
	timestamp := "timestamp"
	lastHash := "lastHash"
	data := []Transaction{}
	dataInByteSlice, _ := json.Marshal(data)
	dataInString := string(dataInByteSlice)
	expectedBlockHash := chainutil.Hash(timestamp + lastHash + dataInString)

	actualBlockHash := HashBlock(timestamp, lastHash, data)

	if expectedBlockHash != actualBlockHash {
		t.Errorf("HashBlock failed\nexpected:\t%v\nactual:\t%v\n", expectedBlockHash, actualBlockHash)
	}
}

func TestSignBlockHash(t *testing.T) {
	wallet := NewWallet("test secret")
	hash := "test hash"
	expectedBlockHash := SignBlockHash(hash, *wallet)
	actualBlockHash := SignBlockHash(hash, *wallet)
	if expectedBlockHash != actualBlockHash {
		t.Errorf("SignBlockHash failed\nexpected:%v\nactual:%v\n", expectedBlockHash, actualBlockHash)
	}
}

func TestCreateBlock(t *testing.T) {
	lastBlock := Genesis()
	data := []Transaction{}
	wallet := NewWallet("test secret")

	actualBlock := CreateBlock(lastBlock, data, *wallet)
	timestamp := actualBlock.timestamp
	hash := HashBlock(timestamp, actualBlock.lastHash, actualBlock.data)
	proposer := chainutil.Key2Str((*wallet).GetPublicKey())
	signature := SignBlockHash(hash, *wallet)

	expectedBlock := Block{
		timestamp,
		lastBlock.hash,
		hash,
		data,
		proposer,
		signature,
		lastBlock.sequenceNum + 1}

	if expectedBlock.timestamp != actualBlock.timestamp {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.timestamp, actualBlock.timestamp)
	}
	if expectedBlock.lastHash != actualBlock.lastHash {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.lastHash, actualBlock.lastHash)
	}
	if expectedBlock.hash != actualBlock.hash {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.hash, actualBlock.hash)
	}
	if len(expectedBlock.data) != len(actualBlock.data) {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.data, actualBlock.data)
	}
	if expectedBlock.proposer != actualBlock.proposer {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.proposer, actualBlock.proposer)
	}
	if expectedBlock.signature != actualBlock.signature {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.signature, actualBlock.signature)
	}
	if expectedBlock.sequenceNum != actualBlock.sequenceNum {
		t.Errorf("CreateBlock failed\nexpected:%v\nactual:%v\n", expectedBlock.sequenceNum, actualBlock.sequenceNum)
	}
}

func TestVerifyBlock(t *testing.T) {
	genesis := Genesis()
	data := []Transaction{}
	wallet := NewWallet("secret")
	block := CreateBlock(genesis, data, *wallet)

	if !VerifyBlock(*block) {
		t.Errorf("VerifyBlock failed!")
	}
}

func TestVerifyProposer(t *testing.T) {
	genesis := Genesis()
	data := []Transaction{}
	wallet1 := NewWallet("secret1")
	wallet2 := NewWallet("secret2")
	block1 := CreateBlock(genesis, data, *wallet1)
	block2 := CreateBlock(genesis, data, *wallet2)
	if !VerifyProposer(*block1, block1.proposer) {
		t.Errorf("VerifyProposer failed!")
	}
	if VerifyProposer(*block2, block1.proposer) {
		t.Errorf("VerifyProposer failed!")
	}
}
