package chain_util2

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"log"
)

// set alias
type PrivateKey = ed25519.PrivateKey
type PublicKey = ed25519.PublicKey

// Hash hashes the data using SHA-256
func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return string(hash[:])
}

// GenKeypair generates keypair with given secret
func GenKeypair(secret string) (PrivateKey, PublicKey) {
	// hash the secret
	hash := Hash(secret)
	seed := []byte(hash)

	// generate Ed25519 keypair from the seed
	privateKey := ed25519.NewKeyFromSeed(seed)
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

// Key2Str converts given privateKey/publicKey to string in hex encoding
func Key2Str(key []byte) string {
	return hex.EncodeToString(key)
}

// Str2Key decodes the string to privateKey/publicKey
func Str2Key(str string) []byte {
	key, err := hex.DecodeString(str)
	if err != nil {
		log.Fatalf("Failed to decode key: %v", err)
	}
	return key
}

// Sign signs given hash with privateKey then returns a signature
func Sign(privateKey PrivateKey, hash string) string {
	signature := ed25519.Sign(privateKey, []byte(hash))
	return string(signature)
}

// Verify verifies the given hash and signature with the given publicKey
func Verify(publicKey PublicKey, hash string, signature string) bool {
	return ed25519.Verify(publicKey, []byte(hash), []byte(signature))
}
