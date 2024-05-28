package pbft

import "testing"

func TestNewRoundChange(t *testing.T) {
	wallet := NewWallet("secret")
	block := NewBlock("timestamp",
		"lastHash",
		"hash",
		[]Transaction{},
		"proposer",
		"signature",
		uint64(0))
	prepare := NewPrepare(*wallet, *block)
	commit := NewCommit(*wallet, *prepare)
	roundChange := NewRoundChange(*wallet, *commit)
	if roundChange.blockHash != commit.blockHash ||
		roundChange.publicKey != commit.publicKey ||
		roundChange.signature != commit.signature {
		t.Error("NewCommit fail")
	}
}

func TestNewRoundChangePool(t *testing.T) {
	roundChangePool := NewRoundChangePool()
	if len(roundChangePool.mapOfList) != 0 {
		t.Error("NewRoundChangePool fail")
	}
}

func TestRoundChangePool_InitRoundChangeForGivenCommit(t *testing.T) {
	roundChangePool := NewRoundChangePool()
	blockHash := "blockHash"
	roundChangePool.InitRoundChangeForGivenCommit(blockHash)
	if len(roundChangePool.mapOfList) != 1 ||
		len(roundChangePool.mapOfList[blockHash]) != 0 {
		t.Error("InitRoundChangeForGivenCommit fail")
	}
}

func TestRoundChangePool_AddRoundChange2Pool(t *testing.T) {
	roundChangePool := NewRoundChangePool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	commit := NewCommit(*wallet, *prepare)
	roundChange := NewRoundChange(*wallet, *commit)
	roundChangePool.InitRoundChangeForGivenCommit(roundChange.blockHash)
	roundChangePool.AddRoundChange2Pool(*roundChange)
	if roundChangePool.mapOfList[roundChange.blockHash][0] != *roundChange {
		t.Error("AddRoundChange2Pool fail")
	}
}

func TestRoundChangePool_RoundChangeExists(t *testing.T) {
	roundChangePool := NewRoundChangePool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	commit := NewCommit(*wallet, *prepare)
	roundChange := NewRoundChange(*wallet, *commit)
	roundChangePool.InitRoundChangeForGivenCommit(roundChange.blockHash)
	if roundChangePool.RoundChangeExists(*roundChange) {
		t.Error("RoundChangeExists fail")
	}
	roundChangePool.AddRoundChange2Pool(*roundChange)
	if !roundChangePool.RoundChangeExists(*roundChange) {
		t.Error("RoundChangeExists fail")
	}
}

func TestRoundChangePool_IsRoundChangeValid(t *testing.T) {
	roundChangePool := NewRoundChangePool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	commit := NewCommit(*wallet, *prepare)
	roundChange := NewRoundChange(*wallet, *commit)
	if !roundChangePool.IsRoundChangeValid(*roundChange) {
		t.Error("IsRoundChangeValid fail")
	}
}
