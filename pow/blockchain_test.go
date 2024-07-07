package pow

import (
	"testing"
	"time"
)

func TestBlockchain_VerifyBlockWithLastBlockInChain(t *testing.T) {
	validators := NewValidators(1)
	bc := NewBlockchain(*validators)
	wallet := NewWallet("Node-1")
	lastHash := bc.chain[len(bc.chain)-1].Hash
	data := make([]Transaction, 0)
	nonce := uint64(1)
	hash := HashBlock(lastHash, data, nonce)
	signature := wallet.Sign(hash)
	block1 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	if !bc.VerifyBlockWithLastBlockInChain(*block1) {
		t.Error("1. VerifyBlockWithLastBlockInChain failed")
	}
	bc.AddBlock(*block1)
	if bc.VerifyBlockWithLastBlockInChain(*block1) {
		t.Error("2. VerifyBlockWithLastBlockInChain failed")
	}
}

func TestBlockchain_BlockConflicts(t *testing.T) {
	validators := NewValidators(1)
	bc := NewBlockchain(*validators)
	wallet := NewWallet("Node-1")
	lastHash := bc.chain[len(bc.chain)-1].Hash
	data := make([]Transaction, 0)
	nonce := uint64(1)
	hash := HashBlock(lastHash, data, nonce)
	signature := wallet.Sign(hash)
	block1 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	block2 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	bc.AddBlock(*block2)
	if !bc.BlockConflicts(*block1) {
		t.Error("1. BlockConflicts failed")
	}
	lastHash = bc.chain[len(bc.chain)-1].Hash
	hash = HashBlock(lastHash, data, nonce)
	signature = wallet.Sign(hash)
	block3 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	if bc.BlockConflicts(*block3) {
		t.Error("2. BlockConflicts failed")
	}
}

func TestBlockchain_Overwritable(t *testing.T) {
	validators := NewValidators(1)
	bc := NewBlockchain(*validators)
	wallet := NewWallet("Node-1")
	lastHash := bc.chain[len(bc.chain)-1].Hash
	data := make([]Transaction, 0)
	nonce := uint64(1)
	hash := HashBlock(lastHash, data, nonce)
	signature := wallet.Sign(hash)
	block1 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	block2 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	bc.AddBlock(*block2)
	if !bc.Overwritable(*block1) {
		t.Error("1. Overwritable failed")
	}
}

func TestBlockchain_Overwritable2(t *testing.T) {
	validators := NewValidators(1)
	bc := NewBlockchain(*validators)
	wallet := NewWallet("Node-1")
	lastHash := bc.chain[len(bc.chain)-1].Hash
	data := make([]Transaction, 0)
	nonce := uint64(1)
	hash := HashBlock(lastHash, data, nonce)
	signature := wallet.Sign(hash)
	block1 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	block2 := wallet.CreateBlock(
		time.Now().String(),
		lastHash,
		hash,
		data,
		wallet.pubKey,
		signature,
		nonce,
	)
	bc.AddBlock(*block1)
	if bc.Overwritable(*block2) {
		t.Error("1. Overwritable failed")
	}
}
