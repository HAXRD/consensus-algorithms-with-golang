package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"testing"
)

func TestNewPrepare(t *testing.T) {
	wallet := NewWallet("secret")
	timestamp := "timestamp"
	lastHash := "lastHash"
	hash := "hash"
	data := []Transaction{}
	proposer := chainutil.Key2Str(wallet.publicKey)
	signature := wallet.Sign(hash)
	sequenceNum := uint64(1)

	block := NewBlock(timestamp, lastHash, hash, data, proposer, signature, sequenceNum)
	prepare := NewPrepare(*wallet, *block)
	if block.hash != prepare.blockHash ||
		chainutil.Key2Str(wallet.publicKey) != prepare.publicKey ||
		block.signature != prepare.signature {
		t.Error("NewPrepare failed")
	}
}

func TestNewPreparePool(t *testing.T) {
	preparePool := NewPreparePool()
	if len(preparePool.mapOfList) != 0 {
		t.Error("NewPreparePool failed")
	}
}

func TestPreparePool_InitPrepareForGivenBlock(t *testing.T) {
	preparePool := NewPreparePool()
	blockHash := "blockHash"
	preparePool.InitPrepareForGivenBlock(blockHash)
	if len(preparePool.mapOfList) != 1 ||
		len(preparePool.mapOfList[blockHash]) != 0 {
		t.Error("InitPrepareForGivenBlock failed")
	}
}

func TestPreparePool_AddPrepare2Pool(t *testing.T) {
	preparePool := NewPreparePool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	preparePool.InitPrepareForGivenBlock(block.hash)
	preparePool.AddPrepare2Pool(*prepare)
	if preparePool.mapOfList[block.hash][0] != *prepare {
		t.Error("AddPrepare2Pool failed")
	}
}

func TestPreparePool_PrepareExists(t *testing.T) {
	preparePool := NewPreparePool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	preparePool.InitPrepareForGivenBlock(block.hash)
	if preparePool.PrepareExists(*prepare) {
		t.Error("PrepareExists failed")
	}
	preparePool.AddPrepare2Pool(*prepare)
	if !preparePool.PrepareExists(*prepare) {
		t.Error("PrepareExists failed")
	}
}

func TestPreparePool_IsPrepareValid(t *testing.T) {
	preparePool := NewPreparePool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	if !preparePool.IsPrepareValid(*prepare) {
		t.Error("IsPrepareValid failed")
	}
}
