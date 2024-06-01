package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
	"log"
	"time"
)

/*
*
Message contains data and timestamp when the tx is created, featured with the following methods:
1. NewMessage
*/

type Message struct {
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

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

/**
TransactionPool temporarily stores txs made by different wallets for each node
It features with the following methods:
1. NewTxPool
2. TxExists
3. AddTx2Pool
4. VerifyTx
5. CleanPool
*/

type TransactionPool struct {
	txs []*Transaction
}

// NewMessage creates a message with given data and timestamp
func NewMessage(data string) *Message {
	return &Message{
		Data:      data,
		Timestamp: time.Now().String(),
	}
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

// NewTxPool creates a tx pool that temporarily stores the txs from
// all available nodes. Txs in pool will be periodically removed
// by matching tx's id.
// TODO: make sure this removing logic works!
func NewTxPool() *TransactionPool {
	return &TransactionPool{
		txs: make([]*Transaction, 0, 2*TX_THRESHOLD),
	}
}

// TxExists checks if a tx exists in the pool or not
func (tp *TransactionPool) TxExists(tx Transaction) bool {
	for _, _tx := range tp.txs {
		if _tx.Id == tx.Id {
			return true
		}
	}
	return false
}

// AddTx2Pool adds a given tx's address to the pool
// returns true if it reaches
func (tp *TransactionPool) AddTx2Pool(tx Transaction) bool {
	tp.txs = append(tp.txs, &tx)
	if len(tp.txs) >= TX_THRESHOLD {
		return true
	}
	return false
}

// VerifyTx checks if a given tx is valid or not
func (tp *TransactionPool) VerifyTx(tx Transaction) bool {
	return tx.VerifyTx()
}

// CleanPool cleans all txs exist in the given block
func (tp *TransactionPool) CleanPool(txs []*Transaction) {
	newTxs := make([]*Transaction, 0, len(tp.txs))
	for _, tx := range txs {
		if !tp.TxExists(*tx) {
			newTxs = append(newTxs, tx)
		}
	}
	tp.txs = newTxs
}