package chainutil

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"log"

	"github.com/google/uuid"
)

type PrivateKey = ed25519.PrivateKey
type PublicKey = ed25519.PublicKey

// GenKeypair generate keypair with given secret
// TODO: consider move this to Wallet
func GenKeypair(secret string) (PrivateKey, PublicKey) {
	// hash the secret using SHA-512
	hash := sha512.Sum512([]byte(secret))
	seed := hash[:32] // Ed25519 seed takes a 32 bytes input

	// generate Ed25519 keypair from the seed
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return privateKey, publicKey
}

// Id return a uuid
func Id() string {
	uuidV1, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil.String()
	}
	return uuidV1.String()
}

// Hash given data using SHA256
func Hash(data string) string {
	hash := sha512.Sum512([]byte(data))
	return string(hash[:])
}

// Sign the given data with privateKey
// TODO: Consider move this to Wallet
func Sign(privateKeyStr string, data string) string {
	privateKey := Str2Key(privateKeyStr, false)
	signature := ed25519.Sign(privateKey, []byte(data))
	return string(signature)
}

// Verify the signed hash (message)
// with publicKey and signature
// TODO: consider move this to Wallet
func Verify(publicKeyStr string, message string, signature string) bool {
	publicKey := Str2Key(publicKeyStr, true)
	hash := Hash(message)
	return ed25519.Verify(publicKey, []byte(hash), []byte(signature))
}

// Key2Str converts given Key to string
func Key2Str(key []byte) string {
	return hex.EncodeToString(key)
}

// Str2Key decodes string to key
func Str2Key(str string, isPublic bool) []byte {
	key, err := hex.DecodeString(str)
	if err != nil {
		log.Fatalf("Failed to decode the given string: %v", err)
	}
	if !isPublic {
		return PrivateKey(key)
	} else {
		return PublicKey(key)
	}
}
