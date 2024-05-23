package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"encoding/json"
	"time"
)

/*
*
1. constructor
2. verify tx
*/

type Message struct {
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

type Transaction struct {
	id        string `json:"id"`
	from      string `json:"from"`
	to        string `json:"to"`
	message   string `json:"message"`
	hash      string `json:"hash"`
	signature string `json:"signature"`
}

func NewMessage(data string) *Message {
	return &Message{
		Data:      data,
		Timestamp: time.Now().String(),
	}
}

func NewTransaction(wallet Wallet, to string, data string) *Transaction {
	messageInJson, _ := json.Marshal(NewMessage(data))
	message := string(messageInJson)
	hash := chainutil.Hash(message)
	signature := wallet.Sign(hash)

	return &Transaction{
		id:        chainutil.Id(),
		from:      chainutil.Key2Str(wallet.GetPublicKey()),
		to:        to,
		message:   message,
		hash:      hash,
		signature: signature,
	}
}

// VerifyTransaction verifies the given transaction with data inside
func VerifyTransaction(transaction Transaction) bool {
	return chainutil.Verify(
		transaction.from,
		transaction.message,
		transaction.signature,
	)
}
