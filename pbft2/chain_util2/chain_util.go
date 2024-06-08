package chain_util2

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"strconv"
)

// set alias
type PrivateKey = ed25519.PrivateKey
type PublicKey = ed25519.PublicKey

// Hash hashes the data using SHA-256
func Hash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// BytesToHex encodes the binary slice into hex code
func BytesToHex(src []byte) string {
	return hex.EncodeToString(src)
}

// HexToBytes decodes the hex code to a binary slice
func HexToBytes(src string) ([]byte, error) {
	res, err := hex.DecodeString(src)
	return res, err
}

// FormatUrl formats url with given host and port
func FormatUrl(host string, port uint64) string {
	return host + ":" + strconv.FormatUint(port, 10)
}

// GenKeypair generates keypair with given secret
func GenKeypair(secret string) (PrivateKey, PublicKey) {
	// hash the secret
	hash := Hash(secret)

	// generate Ed25519 keypair from the seed
	privateKey := ed25519.NewKeyFromSeed(hash)
	publicKey := privateKey.Public().(PublicKey)

	return privateKey, publicKey
}

// Id returns a uuid
func Id() string {
	uuidV1, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil.String()
	}
	return uuidV1.String()
}

// Sign signs given hash with privateKey then returns a signature
func Sign(privateKey PrivateKey, hash []byte) []byte {
	signature := ed25519.Sign(privateKey, hash)
	return signature
}

// Verify verifies the given hash and signature with the given publicKey
func Verify(publicKey PublicKey, hash []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, hash, signature)
}
