package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
	"crypto/ed25519"
	"fmt"
	"log"
	"time"
)

/**
Wallet features the following methods:
1. NewWallet
2. PrintWallet
3. Sign
4. Verify
5. CreateTx
6. CreateBlock
*/

// set alias
type PrivateKey = ed25519.PrivateKey
type PublicKey = ed25519.PublicKey

// Wallet contains the private-public keypair
type Wallet struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// NewWallet creates a new wallet by generating a keypair with given secret
func NewWallet(secret string) *Wallet {
	privateKey, publicKey := chain_util.GenKeypair(secret)
	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// PrintWallet prints wallet's publicKey
func (w *Wallet) PrintWallet() {
	fmt.Printf("Wallet - public key: %s\n", chain_util.BytesToHex(w.publicKey)[:6])
}

// Sign uses wallet's privateKey to sign a given hash and returns a signature
func (w *Wallet) Sign(hash []byte) []byte {
	signature := chain_util.Sign(w.privateKey, hash)
	return signature
}

// CreateTx creates a tx with given data
func (w *Wallet) CreateTx(data string) *Transaction {
	return NewTx(*w, data)
}

// CreateBlock creates a block with lastBlock and provided data
func (w *Wallet) CreateBlock(lastBlock Block, data []Transaction) *Block {
	timestamp := time.Now().String()
	lastHash := lastBlock.Hash
	nonce := lastBlock.Nonce + 1
	// hash block with timestamp, lastBlock's hash, marshalled data and current nonce
	hash := HashBlock(timestamp, lastHash, data, nonce)
	// sign the hash
	signature := w.Sign(hash)
	block := NewBlock(
		timestamp,
		lastHash,
		hash,
		data,
		w.publicKey,
		signature,
		nonce,
		nil, nil, nil, nil,
	)
	log.Printf("Created block [%s]\n", chain_util.BytesToHex(block.Hash)[:6])
	return block
}

// CreateMsg creates a message for PBFT phase transition
func (w *Wallet) CreateMsg(msgType string, blockHash []byte) *Message {
	return NewMsg(
		msgType,
		blockHash,
		w.publicKey,
		w.Sign(blockHash),
	)
}
