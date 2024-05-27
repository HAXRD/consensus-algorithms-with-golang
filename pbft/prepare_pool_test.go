package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"testing"
)

func TestNewPrepare(t *testing.T) {
	blockHash := "blockHash"
	publicKey := "publicKey"
	signature := "signature"
	prepare := NewPrepare(blockHash, publicKey, signature)

	if blockHash != prepare.blockHash ||
		publicKey != prepare.publicKey ||
		signature != prepare.signature {
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
	blockHash := "blockHash"
	publicKey := "publicKey"
	signature := "signature"
	prepare := NewPrepare(blockHash, publicKey, signature)
	preparePool.InitPrepareForGivenBlock(blockHash)
	preparePool.AddPrepare2Pool(*prepare)
	if preparePool.mapOfList[blockHash][0] != *prepare {
		t.Error("AddPrepare2Pool failed")
	}
}

func TestPreparePool_PrepareExists(t *testing.T) {
	preparePool := NewPreparePool()
	blockHash := "blockHash"
	publicKey := "publicKey"
	signature := "signature"
	prepare := NewPrepare(blockHash, publicKey, signature)
	preparePool.InitPrepareForGivenBlock(blockHash)
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
	blockHash := "blockHash"
	privateKey, publicKey := chainutil.GenKeypair("secret")
	privateKeyStr := chainutil.Key2Str(privateKey)
	publicKeyStr := chainutil.Key2Str(publicKey)
	signature := chainutil.Sign(privateKeyStr, blockHash)
	prepare := NewPrepare(blockHash, publicKeyStr, signature)
	if !preparePool.IsPrepareValid(*prepare) {
		t.Error("IsPrepareValid failed")
	}
}
