package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"crypto/ed25519"
	"fmt"
)

/**
1. constructor
2. print Wallet
3. sign
4. verify
5. create tx
6. get publicKey
*/

// Wallet struct that contains a keypair
type Wallet struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewWallet create a new wallet with given secret
func NewWallet(secret string) *Wallet {
	privateKey, publicKey := chainutil.GenKeypair(secret)
	return &Wallet{privateKey: privateKey, publicKey: publicKey}
}

// PrintWallet prints the publicKey of the given wallet
func (w *Wallet) PrintWallet() {
	fmt.Printf("Wallet - public key: %s\n", w.publicKey)
}

// Sign the given data
func (w *Wallet) Sign(data string) string {
	signature := chainutil.Sign(chainutil.Key2Str(w.privateKey), data)
	return signature
}

// Verify the hash with given wallet and signature
func (w *Wallet) Verify(hash string, signature string) bool {
	return chainutil.Verify(chainutil.Key2Str(w.publicKey), hash, signature)
}

// CreateTransaction creates a new transaction
func (w *Wallet) CreateTransaction(to string, data string) Transaction {
	return *NewTransaction(*w, to, data)
}

// GetPublicKey gets the publicKey for given wallet
func (w *Wallet) GetPublicKey() ed25519.PublicKey {
	return w.publicKey
}
