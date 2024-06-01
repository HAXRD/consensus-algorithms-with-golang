package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"crypto/ed25519"
	"fmt"
)

/**
Wallet features the following methods:
1. NewWallet
2. PrintWallet
3. Sign
4. Verify
5. CreateTx
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
	privateKey, publicKey := chain_util2.GenKeypair(secret)
	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// PrintWallet prints wallet's publicKey
func (w *Wallet) PrintWallet() {
	fmt.Printf("Wallet - public key: %s\n", chain_util2.Key2Str(w.publicKey))
}

// Sign uses wallet's privateKey to sign a given hash and returns a signature
func (w *Wallet) Sign(hash string) string {
	signature := chain_util2.Sign(w.privateKey, hash)
	return signature
}

// Verify verifies given hash and signature is matched with message
// TODO: suspect this method will never be used. Remove it if so
//func (w *Wallet) Verify(hash string, signature string) bool {
//	return chain_util2.Verify(w.publicKey, hash, signature)
//}

// CreateTx creates a tx with given data
func (w *Wallet) CreateTx(data string) Transaction {
	return *NewTx(*w, data)
}
