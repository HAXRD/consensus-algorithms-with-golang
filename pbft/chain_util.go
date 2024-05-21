package pbft

import (
	"crypto/ed25519"
	"crypto/sha512"

	"github.com/google/uuid"
)

type PrivateKey = ed25519.PrivateKey
type PublicKey = ed25519.PublicKey

// GenKeypair generate keypair with given secret
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
func Sign(privateKey PrivateKey, data string) string {
	signature := ed25519.Sign(privateKey, []byte(data))
	return string(signature)
}

// Verify the signed hash (message)
// with publicKey and signature
func Verify(publicKey PublicKey, message string, signature string) bool {
	return ed25519.Verify(publicKey, []byte(message), []byte(signature))
}
