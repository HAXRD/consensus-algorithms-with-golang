package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"testing"
)

// test NewMessage function
func TestNewMessage(t *testing.T) {
	data := "This a test data"
	message := NewMessage(data)

	if message.Data != data {
		t.Errorf("NewMessage failed, expected %s, got %s", data, message.Data)
	}
}

// test NewTransaction function
func TestNewTransaction(t *testing.T) {
	w1 := NewWallet("wallet1")
	w2 := NewWallet("wallet2")
	to := chainutil.Key2Str(w2.GetPublicKey())
	data := "w1 -> w2"

	transaction := NewTransaction(*w1, to, data)

	if chainutil.Key2Str(w1.GetPublicKey()) != transaction.from {
		t.Errorf(
			"NewTransaction failed, expected from %s, got from %s",
			chainutil.Key2Str(w1.GetPublicKey()),
			transaction.from)
	}
	if chainutil.Key2Str(w2.GetPublicKey()) != transaction.to {
		t.Errorf(
			"NewTransaction failed, expected to %s, got to %s",
			chainutil.Key2Str(w2.GetPublicKey()),
			transaction.to)
	}
}

// test VerifyTransaction function
func TestVerifyTransaction(t *testing.T) {
	w1 := NewWallet("wallet1")
	w2 := NewWallet("wallet2")
	data := "w1 -> w2"
	tx := NewTransaction(*w1, chainutil.Key2Str(w2.GetPublicKey()), data)

	if !VerifyTransaction(*tx) {
		t.Errorf("VerifyTransaction failed, expected true, got false")
	}
}
