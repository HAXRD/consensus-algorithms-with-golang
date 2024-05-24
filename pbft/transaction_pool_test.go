package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"testing"
)

func TestNewTransactionPool(t *testing.T) {
	txPool := NewTransactionPool()
	if len(txPool.transactions) != 0 {
		t.Errorf("NewTransactionPool returned wrong number of transactions, expected 0, got %v", len(txPool.transactions))
	}
	wallet1 := NewWallet("secret1")
	wallet2 := NewWallet("secret2")
	to := chainutil.Key2Str(wallet2.GetPublicKey())
	data := "this is a test message"
	tx := NewTransaction(*wallet1, to, data)
	txPool.AddTransaction(*tx)
	if len(txPool.transactions) != 1 {
		t.Errorf("NewTransactionPool returned wrong number of transactions, expected 1, got %v", len(txPool.transactions))
	}
}

func TestTransactionPool_TransactionExists(t *testing.T) {
	txPool := NewTransactionPool()
	wallet1 := NewWallet("secret1")
	wallet2 := NewWallet("secret2")
	to := chainutil.Key2Str(wallet2.GetPublicKey())
	data := "this is a test message"
	tx := NewTransaction(*wallet1, to, data)

	if txPool.TransactionExists(*tx) {
		t.Errorf("TransactionExists returned wrong transaction")
	}

	txPool.AddTransaction(*tx)

	if !txPool.TransactionExists(*tx) {
		t.Errorf("TransactionExists returned wrong transaction")
	}
}

func TestTransactionPool_AddTransaction(t *testing.T) {
	txPool := NewTransactionPool()
	wallet1 := NewWallet("secret1")
	wallet2 := NewWallet("secret2")
	to := chainutil.Key2Str(wallet2.GetPublicKey())
	data := "this is a test message"
	tx := NewTransaction(*wallet1, to, data)

	if txPool.AddTransaction(*tx) {
		t.Errorf("AddTransaction returned wrong transaction")
	}
	if txPool.AddTransaction(*tx) {
		t.Errorf("AddTransaction returned wrong transaction")
	}
	if txPool.AddTransaction(*tx) {
		t.Errorf("AddTransaction returned wrong transaction")
	}
	if txPool.AddTransaction(*tx) {
		t.Errorf("AddTransaction returned wrong transaction")
	}
	if !txPool.AddTransaction(*tx) {
		t.Errorf("AddTransaction returned wrong transaction")
	}
}

func TestTransactionPool_VerifyTransaction(t *testing.T) {
	txPool := NewTransactionPool()
	wallet1 := NewWallet("secret1")
	wallet2 := NewWallet("secret2")
	to := chainutil.Key2Str(wallet2.GetPublicKey())
	data := "this is a test message"
	tx := NewTransaction(*wallet1, to, data)
	if !txPool.VerifyTransaction(*tx) {
		t.Errorf("VerifyTransaction returned wrong transaction")
	}
}
