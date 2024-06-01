package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
	"log"
	"time"
)

/**
Transaction is created by a wallet, featured with the following methods:
1. NewTx
2. VerifyTx
*/

type Transaction struct {
	Id        string `json:"id"`
	From      string `json:"from"`
	Message   string `json:"message"`
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

// NewTx create a tx with a wallet
func NewTx(w Wallet, data string) *Transaction {
	msg := NewMessage(data)
	msgStr, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Tx's msg json marshal err, %v\n", err)
	}
	hash := chain_util2.Hash(string(msgStr))
	signature := w.Sign(hash)

	return &Transaction{
		Id:        chain_util2.Id(),
		From:      chain_util2.Key2Str(w.publicKey),
		Message:   string(msgStr),
		Hash:      hash,
		Signature: signature,
	}
}

// VerifyTx verifies a given tx with tx's msg->hash and hash->signature
func (tx *Transaction) VerifyTx() bool {
	// verify msg->hash
	if tx.Hash != chain_util2.Hash(tx.Message) {
		return false
	}
	// verify hash->signature
	publicKey := chain_util2.Str2Key(tx.From)
	if !chain_util2.Verify(publicKey, tx.Hash, tx.Signature) {
		return false
	}
	return true
}

/*
*
Message contains data and timestamp when the tx is created, featured with the following methods:
1. NewMessage
*/

type Message struct {
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

// NewMessage creates a message with given data and timestamp
func NewMessage(data string) *Message {
	return &Message{
		Data:      data,
		Timestamp: time.Now().String(),
	}
}
