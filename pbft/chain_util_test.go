package pbft

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
)

// test the GenKeypair function
func TestGenKeypair(t *testing.T) {
	secret := "test-secret-string"
	hash := sha512.Sum512([]byte(secret))
	seed := hash[:32]
	expectedPrivateKey := ed25519.NewKeyFromSeed(seed)
	expectedPublicKey := expectedPrivateKey.Public().(ed25519.PublicKey)

	privateKey, publicKey := GenKeypair(secret)

	if hex.EncodeToString(expectedPrivateKey) != hex.EncodeToString(privateKey) {
		t.Errorf("Expected private key, %s, got, %s", expectedPrivateKey, privateKey)
	}

	if hex.EncodeToString(expectedPublicKey) != hex.EncodeToString(publicKey) {
		t.Errorf("Expected public key, %s, got, %s", expectedPublicKey, publicKey)
	}
}

// test GenKeypair if the generated keypair is valid
func TestGenerateKeypairIsValid(t *testing.T) {
	secret := "test-secret-string"
	privateKey, publicKey := GenKeypair(secret)

	message := []byte("A test message")
	signature := ed25519.Sign(privateKey, message)
	if !ed25519.Verify(publicKey, message, signature) {
		t.Errorf("Failed to verify signature with generated key pair!")
	}
}

// test the Id function
func TestId(t *testing.T) {
	// test uniqueness
	id1 := Id()
	id2 := Id()
	if id1 == id2 {
		t.Errorf("Expected unique UUIDs, but got the same, %s", id1)
	}

	// test validity
	if _, err := uuid.Parse(id1); err != nil {
		t.Errorf("Generate UUID is not valid, %s", id1)
	}
	if _, err := uuid.Parse(id2); err != nil {
		t.Errorf("Generate UUID is not valid, %s", id2)
	}
}

// test the Hash function
func TestHash(t *testing.T) {
	data := "Test data"
	hash := sha512.Sum512([]byte(data))
	expected := string(hash[:])

	actual := Hash(data)

	if expected != actual {
		t.Errorf("Expected, %s\nGot, %s", hex.EncodeToString([]byte(expected)), hex.EncodeToString([]byte(actual)))
	}
}

// test the Sign function
func TestSign(t *testing.T) {
	secret := "test-secret"
	privateKey, _ := GenKeypair(secret)
	message := "A test message"
	expected := string(ed25519.Sign(privateKey, []byte(message)))

	actual := Sign(privateKey, message)

	if expected != actual {
		t.Errorf("Expected, %s\nGot, %s", expected, actual)
	}
}

// test the Verify function
func TestVerify(t *testing.T) {
	secret := "test secret"
	message := "This is a test message"
	privateKey, publicKey := GenKeypair(secret)
	signature := string(ed25519.Sign(privateKey, []byte(message)))

	if !Verify(publicKey, message, signature) {
		t.Errorf("Failed to verify signature with generated key pair!")
	}
}
