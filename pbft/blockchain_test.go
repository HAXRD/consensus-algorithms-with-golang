package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"fmt"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockchain()
	if len(blockchain.validatorList) != NUM_OF_NODES ||
		len(blockchain.chain) != 1 {
		t.Error("NewBlockchain fail")
	}
}

func TestBlockchain_CreateBlock(t *testing.T) {
	blockchain := NewBlockchain()
	wallet := NewWallet("secret")
	to := chainutil.Key2Str(NewWallet("to").privateKey)
	data := "w1 -> w2"
	txs := make([]Transaction, 0)
	txs = append(txs, *NewTransaction(*wallet, to, data))
	block := blockchain.CreateBlock(txs, *wallet)
	if block.lastHash != blockchain.chain[0].hash {
		t.Error("CreateBlock fail")
	}
}

func TestBlockchain_GetProposer(t *testing.T) {
	blockchain := NewBlockchain()
	fmt.Println(blockchain.chain[len(blockchain.chain)-1].hash[0])
	index := blockchain.chain[len(blockchain.chain)-1].hash[0] % NUM_OF_NODES
	if blockchain.GetProposer() != blockchain.validatorList[index] {
		t.Error("GetProposer fail")
	}
}

func TestBlockchain_IsBlockValid(t *testing.T) {
	blockchain := NewBlockchain()
	wallet := NewWallet("secret")
	to := chainutil.Key2Str(NewWallet("to").privateKey)
	data := "w1 -> w2"
	txs := make([]Transaction, 0)
	txs = append(txs, *NewTransaction(*wallet, to, data))
	block := blockchain.CreateBlock(txs, *wallet)
	if !blockchain.IsBlockValid(*block) {
		t.Error("IsBlockValid fail")
	}
}

func TestBlockchain_AddUpdatedBlock2Chain(t *testing.T) {
	blockchain := NewBlockchain()
	wallet := NewWallet("secret")
	to := chainutil.Key2Str(NewWallet("to").privateKey)
	data := "w1 -> w2"
	txs := make([]Transaction, 0)
	txs = append(txs, *NewTransaction(*wallet, to, data))
	block := blockchain.CreateBlock(txs, *wallet)
	blockPool := NewBlockPool()
	blockPool.AddBlock2Pool(*block)
	preparePool := NewPreparePool()
	prepare := NewPrepare(*wallet, *block)
	preparePool.AddPrepare2Pool(*prepare)
	commitPool := NewCommitPool()
	commit := NewCommit(*wallet, *prepare)
	commitPool.AddCommit2Pool(*commit)
	blockchain.AddUpdatedBlock2Chain(block.hash, *blockPool, *preparePool, *commitPool)
	if len(blockchain.chain) != 2 {
		t.Error("AddUpdatedBlock2Chain fail")
	}
}
