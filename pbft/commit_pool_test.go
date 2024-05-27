package pbft

import (
	"testing"
)

func TestNewCommit(t *testing.T) {
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
	if commit.blockHash != prepare.blockHash ||
		commit.publicKey != prepare.publicKey ||
		commit.signature != prepare.signature {
		t.Error("NewCommit fail")
	}
}

func TestNewCommitPool(t *testing.T) {
	commitPool := NewCommitPool()
	if len(commitPool.mapOfList) != 0 {
		t.Error("NewCommitPool fail")
	}
}

func TestCommitPool_InitCommitForGivenPrepare(t *testing.T) {
	commitPool := NewCommitPool()
	blockHash := "blockHash"
	commitPool.InitCommitForGivenPrepare(blockHash)
	if len(commitPool.mapOfList) != 1 ||
		len(commitPool.mapOfList[blockHash]) != 0 {
		t.Error("InitCommitForGivenPrepare fail")
	}
}

func TestCommitPool_AddCommit2Pool(t *testing.T) {
	commitPool := NewCommitPool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	commit := NewCommit(*wallet, *prepare)
	commitPool.InitCommitForGivenPrepare(commit.blockHash)
	commitPool.AddCommit2Pool(*commit)
	if commitPool.mapOfList[commit.blockHash][0] != *commit {
		t.Error("AddCommit2Pool fail")
	}
}

func TestCommitPool_CommitExists(t *testing.T) {
	commitPool := NewCommitPool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	commit := NewCommit(*wallet, *prepare)
	commitPool.InitCommitForGivenPrepare(commit.blockHash)
	if commitPool.CommitExists(*commit) {
		t.Error("CommitExists fail")
	}
	commitPool.AddCommit2Pool(*commit)
	if !commitPool.CommitExists(*commit) {
		t.Error("CommitExists fail")
	}
}

func TestCommitPool_IsCommitValid(t *testing.T) {
	commpool := NewCommitPool()
	wallet := NewWallet("secret")
	block := Genesis()
	prepare := NewPrepare(*wallet, block)
	commit := NewCommit(*wallet, *prepare)
	if !commpool.IsCommitValid(*commit) {
		t.Error("IsCommitValid fail")
	}
}
