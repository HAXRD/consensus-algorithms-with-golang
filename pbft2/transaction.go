package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
	"log"
	"time"
)

/*
*
Event contains data and timestamp when the tx is created,
featured with the following methods:
1. NewEvent
*/

type Event struct {
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
	From      []byte `json:"from"`
	Event     Event  `json:"event"`
	Hash      []byte `json:"hash"`
	Signature []byte `json:"signature"`
	MsgType   string `json:"msgType"`
}

// MarshalJSON is the custom Marshal function for Transaction
func (tx Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		From      string `json:"from"`
		Hash      string `json:"hash"`
		Signature string `json:"signature"`
		Alias
	}{
		From:      chain_util2.BytesToHex(tx.From),
		Hash:      chain_util2.BytesToHex(tx.Hash),
		Signature: chain_util2.BytesToHex(tx.Signature),
		Alias:     Alias(tx),
	})
}

// UnmarshalJSON is the custom Unmarshal function for Transaction
func (tx *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		From      string `json:"from"`
		Hash      string `json:"hash"`
		Signature string `json:"signature"`
		*Alias
	}{
		Alias: (*Alias)(tx),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	from, err := chain_util2.HexToBytes(aux.From)
	if err != nil {
		return err
	}
	hash, err := chain_util2.HexToBytes(aux.Hash)
	if err != nil {
		return err
	}
	signature, err := chain_util2.HexToBytes(aux.Signature)
	if err != nil {
		return err
	}
	tx.From = from
	tx.Hash = hash
	tx.Signature = signature
	return nil
}

/**
TransactionPool temporarily stores pool made by different wallets for each node.
It features the following methods:
1. NewTxPool
2. TxExists
3. AddTx2Pool
4. VerifyTx
5. CleanPool
*/

type TransactionPool struct {
	pool []Transaction
}

// NewEvent creates a message with given data and timestamp
func NewEvent(data string) *Event {
	return &Event{
		Data:      data,
		Timestamp: time.Now().String(),
	}
}

// NewTx create a tx with a wallet
func NewTx(w Wallet, data string) *Transaction {
	event := NewEvent(data)
	eventStr, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Tx's event json marshal err, %v\n", err)
	}
	hash := chain_util2.Hash(string(eventStr))
	signature := w.Sign(hash)

	return &Transaction{
		Id:        chain_util2.Id(),
		From:      w.publicKey,
		Event:     *event,
		Hash:      hash,
		Signature: signature,
		MsgType:   MsgTx,
	}
}

// VerifyTx verifies a given tx with tx's msg->hash and hash->signature
func (tx *Transaction) VerifyTx() bool {
	eventStr, err := json.Marshal(tx.Event)
	if err != nil {
		log.Fatalf("Tx's event json marshal err, %v\n", err)
	}
	return tx.MsgType == MsgTx && // verify msgType
		chain_util2.BytesToHex(tx.Hash) == chain_util2.BytesToHex(chain_util2.Hash(string(eventStr))) && // verify msg->hash
		chain_util2.Verify(tx.From, tx.Hash, tx.Signature) // verify hash->signature
}

// NewTxPool creates a tx pool that temporarily stores the pool from
// all available nodes. Txs in pool will be periodically removed
// by matching tx's id.
// TODO: make sure this removing logic works!
func NewTxPool() *TransactionPool {
	return &TransactionPool{
		pool: make([]Transaction, 0, 2*TX_THRESHOLD),
	}
}

// TxExists checks if a tx exists in the pool or not
func (tp *TransactionPool) TxExists(tx Transaction) bool {
	for _, _tx := range tp.pool {
		if _tx.Id == tx.Id {
			return true
		}
	}
	return false
}

// AddTx2Pool adds a given tx's address to the pool
// returns true if it reaches
func (tp *TransactionPool) AddTx2Pool(tx Transaction) bool {
	tp.pool = append(tp.pool, tx)
	log.Printf("Tx [%s] added to tx pool\n", chain_util2.BytesToHex(tx.Hash)[:5])
	if len(tp.pool) >= TX_THRESHOLD {
		return true
	}
	return false
}

// VerifyTx checks if a given tx is valid or not
func (tp *TransactionPool) VerifyTx(tx Transaction) bool {
	return tx.VerifyTx()
}

// CleanPool cleans all pool exist in the given block.
// Returns true if any txs have been removed;
// Returns false otherwise.
func (tp *TransactionPool) CleanPool(txs []Transaction) bool {
	if txs == nil || len(txs) == 0 {
		return false
	}
	newTxs := make([]Transaction, 0, len(tp.pool))
	for _, tx := range txs {
		if !tp.TxExists(tx) {
			newTxs = append(newTxs, tx)
		}
	}
	if len(newTxs) == len(tp.pool) {
		return false
	} else {
		tp.pool = newTxs
		return true
	}
}
