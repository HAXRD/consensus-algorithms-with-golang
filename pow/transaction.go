package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
	"encoding/json"
	"log"
	"time"
)

type Event struct {
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

func NewEvent(data string) *Event {
	return &Event{
		Data:      data,
		Timestamp: time.Now().String(),
	}
}

type Transaction struct {
	Id        string `json:"id"`
	From      []byte `json:"from"`
	Event     Event  `json:"event"`
	Hash      []byte `json:"hash"`
	Signature []byte `json:"signature"`
}

func (tx Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		From      string `json:"from"`
		Hash      string `json:"hash"`
		Signature string `json:"signature"`
		Alias
	}{
		From:      pow_util.Byte2Hex(tx.From),
		Hash:      pow_util.Byte2Hex(tx.Hash),
		Signature: pow_util.Byte2Hex(tx.Signature),
		Alias:     Alias(tx),
	})
}

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
	from, err := pow_util.Hex2Byte(aux.From)
	if err != nil {
		return err
	}
	hash, err := pow_util.Hex2Byte(aux.Hash)
	if err != nil {
		return err
	}
	signature, err := pow_util.Hex2Byte(aux.Signature)
	if err != nil {
		return err
	}
	tx.From = from
	tx.Hash = hash
	tx.Signature = signature
	return nil
}

func NewTx(w Wallet, data string) *Transaction {
	event := NewEvent(data)
	eventStr, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	hash := pow_util.Hash(string(eventStr))
	signature := w.Sign(hash)

	return &Transaction{
		Id:        pow_util.Id(),
		From:      w.pubKey,
		Event:     *event,
		Hash:      hash,
		Signature: signature,
	}
}

// VerifyTx verifies the tx with tx' event->hash and hash->signature
func (tx *Transaction) VerifyTx() bool {
	eventStr, err := json.Marshal(tx.Event)
	if err != nil {
		log.Fatalf("Tx's event json marshal err: %v\n", err)
	}
	return pow_util.Byte2Hex(tx.Hash) == pow_util.Byte2Hex(pow_util.Hash(string(eventStr))) && // event->hash
		pow_util.Verify(tx.From, tx.Hash, tx.Signature) // hash->signature
}

type TxPool struct {
	waiting   map[string]Transaction
	committed map[string]Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		waiting:   make(map[string]Transaction),
		committed: make(map[string]Transaction),
	}
}

func (txp *TxPool) TxExists(tx Transaction) bool {
	var ok bool
	// committed
	if _, ok = txp.committed[tx.Id]; ok {
		return true
	}
	// waiting
	if _, ok = txp.waiting[tx.Id]; ok {
		return true
	}
	return false
}

func (tp *TxPool) VerifyTx(tx Transaction) bool {
	return tx.VerifyTx()
}

func (txp *TxPool) AddTx2Pool(tx Transaction) {
	// skip if exists
	if txp.TxExists(tx) {
		return
	}
	txp.waiting[tx.Id] = tx
	log.Printf("Tx [%s] added to tx pool\n", pow_util.Byte2Hex(tx.Hash)[:6])
}

func (txp *TxPool) UpdateCommitted(block Block, blocksThatWereOverwritten []Block) {
	// remove txs from committed && waiting
	var ok bool
	if blocksThatWereOverwritten == nil || len(blocksThatWereOverwritten) == 0 {
		for _, b := range blocksThatWereOverwritten {
			for _, tx := range b.Data {
				if _, ok = txp.committed[tx.Id]; ok {
					delete(txp.committed, tx.Id)
				}
			}
		}
	}
	// remove txs in waiting (if exists)
	// add txs to committed
	for _, tx := range block.Data {
		if _, ok = txp.waiting[tx.Id]; ok {
			delete(txp.waiting, tx.Id)
		}
		txp.committed[tx.Id] = tx
	}
}
