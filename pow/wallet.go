package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
	"crypto/ed25519"
)

type PublicKey = ed25519.PublicKey
type PrivateKey = ed25519.PrivateKey

type Wallet struct {
	priKey PrivateKey
	pubKey PublicKey
}

func NewWallet(secret string) *Wallet {
	priKey, pubKey := pow_util.GenKeyPair(secret)
	return &Wallet{
		priKey: priKey,
		pubKey: pubKey,
	}
}

func (w *Wallet) Sign(data []byte) []byte {
	return pow_util.Sign(w.priKey, data)
}

func (w *Wallet) CreateTx(data string) *Transaction {
	return NewTx(*w, data)
}

func (w *Wallet) CreateBlock(
	timestamp string,
	lastHash []byte,
	hash []byte,
	data []Transaction,
	proposer PublicKey,
	signature []byte,
	nonce uint64,
) *Block {
	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
		Proposer:  proposer,
		Signature: signature,
		Nonce:     nonce,
	}
}
