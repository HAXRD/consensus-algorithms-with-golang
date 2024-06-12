package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
	"encoding/json"
	"testing"
)

func TestNewEvent(t *testing.T) {
	data := "test"
	event1 := NewEvent(data)
	event2 := NewEvent(data)
	if event1 == event2 {
		t.Errorf("event1 and event2 should be different")
	}
	if event1.Data != event2.Data {
		t.Errorf("event1.Data and event2.Data should be the same")
	}

	// marshall & unmarshal
	event := NewEvent(data)
	eventStr, err := json.Marshal(event)
	if err != nil {
		t.Error(err)
	}
	var event3 Event
	err = json.Unmarshal(eventStr, &event3)
	if err != nil {
		t.Error(err)
	}
	if *event != event3 {
		t.Errorf("event and event3 should be the same")
	}
}

func TestNewTx(t *testing.T) {
	w := NewWallet("test")
	data := "test"
	tx1 := NewTx(*w, data)
	tx2 := NewTx(*w, data)
	if tx1 == tx2 {
		t.Errorf("tx1 and tx2 should be different")
	}
	if tx1.Id == tx2.Id ||
		chain_util.BytesToHex(tx1.From) != chain_util.BytesToHex(tx2.From) ||
		tx1.Event == tx2.Event ||
		chain_util.BytesToHex(tx1.Hash) == chain_util.BytesToHex(tx2.Hash) ||
		chain_util.BytesToHex(tx1.Signature) == chain_util.BytesToHex(tx2.Signature) ||
		tx1.MsgType != tx2.MsgType ||
		tx1.MsgType != MsgTx {
		t.Errorf("tx1 and tx2 should be the different")
	}

	// marshal & unmarshal
	tx3 := NewTx(*w, data)
	tx3Str, err := json.Marshal(tx3)
	if err != nil {
		t.Error(err)
	}
	var tx4 Transaction
	err = json.Unmarshal(tx3Str, &tx4)
	if err != nil {
		t.Error(err)
	}
	if tx3.Id != tx4.Id ||
		chain_util.BytesToHex(tx3.From) != chain_util.BytesToHex(tx4.From) ||
		tx3.Event != tx4.Event ||
		chain_util.BytesToHex(tx3.Hash) != chain_util.BytesToHex(tx4.Hash) ||
		chain_util.BytesToHex(tx3.Signature) != chain_util.BytesToHex(tx4.Signature) ||
		tx3.MsgType != tx4.MsgType ||
		tx3.MsgType != MsgTx {
		t.Errorf("tx3 and tx4 should be the same,\ntx3: %+v,\ntx4: %+v\n", *tx3, tx4)
	}
}

func TestTransaction_VerifyTx(t *testing.T) {
	data := "data"
	w := NewWallet("test")
	tx := NewTx(*w, data)
	if !tx.VerifyTx() {
		t.Errorf("VerifyTx should be true")
	}
	tx.Event.Data = "data2"
	if tx.VerifyTx() {
		t.Errorf("VerifyTx should be false")
	}
}

func TestTransactionPool_AddTx2Pool(t *testing.T) {
	data := "data"
	w := NewWallet("test")
	tp := NewTxPool()
	for i := range TX_THRESHOLD {
		tx := NewTx(*w, data)
		poolCopy, _ := tp.AddTx2Pool(*tx)
		if i+1 < TX_THRESHOLD && (poolCopy != nil || len(tp.inProgress) != 0) {
			t.Errorf("AddTx2Pool should return false")
		} else if i+1 >= TX_THRESHOLD && (poolCopy == nil || len(tp.inProgress) == 0) {
			t.Errorf("AddTx2Pool should return true")
		}
	}
}

func TestTransactionPool_TxExists(t *testing.T) {
	data := "data"
	w := NewWallet("test")
	tx1 := NewTx(*w, data)
	tx2 := NewTx(*w, data)
	tx3 := NewTx(*w, data)
	tp := NewTxPool()
	tp.AddTx2Pool(*tx1)
	tp.AddTx2Pool(*tx2)
	if !tp.TxExists(*tx2) {
		t.Errorf("TxExists failed")
	}
	if tp.TxExists(*tx3) {
		t.Errorf("TxExists failed")
	}
	tp.AddTx2Pool(*tx3)
	if !tp.TxExists(*tx3) {
		t.Errorf("TxExists failed")
	}
}

func TestTransactionPool_TransferInProgressToCommitted(t *testing.T) {
	data := "data"
	w := NewWallet("test")
	tx1 := NewTx(*w, data)
	tx2 := NewTx(*w, data)
	tx3 := NewTx(*w, data)
	tp := NewTxPool()
	var returnedTxs []Transaction
	returnedTxs, _ = tp.AddTx2Pool(*tx1)
	if returnedTxs != nil {
		t.Errorf("AddTx2Pool should return nil")
	}
	returnedTxs, _ = tp.AddTx2Pool(*tx2)
	if returnedTxs != nil {
		t.Errorf("AddTx2Pool should return nil")
	}
	if len(tp.inProgress) != 0 {
		t.Errorf("InProgress should be empty")
	}
	if len(tp.committed) != 0 {
		t.Errorf("Committed should be empty")
	}
	returnedTxs, _ = tp.AddTx2Pool(*tx3)
	if returnedTxs == nil {
		t.Errorf("AddTx2Pool should return txs")
	}
	if len(tp.inProgress) == 0 {
		t.Errorf("InProgress should not be empty")
	}
	if len(tp.committed) != 0 {
		t.Errorf("Committed should be empty")
	}
	success := tp.TransferInProgressToCommitted(returnedTxs)
	if !success || len(tp.inProgress) != 0 || len(tp.committed) == 0 {
		t.Errorf("TransferInProgressToCommitted failed")
	}
}
