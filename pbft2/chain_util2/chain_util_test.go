package chain_util2

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	data := "test data"

	hash := sha256.Sum256([]byte(data))
	expected := hash[:]
	actual := Hash(data)

	if BytesToHex(expected) != BytesToHex(actual) {
		t.Errorf("Hash failed, expected %s, actual %s\n", expected, actual)
	}
}

func TestBytesToHex(t *testing.T) {
	data := "test data"
	hash := Hash(data)
	expected := hex.EncodeToString(hash)
	actual := BytesToHex(hash)
	if expected != actual {
		t.Errorf("BytesToHex  failed, expected %s, actual %s\n", expected, actual)
	}
}

func TestHexToBytes(t *testing.T) {
	data := "test data"
	hash := Hash(data)
	hashHex := BytesToHex(hash)
	expected, err := hex.DecodeString(hashHex)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, hash) {
		t.Errorf("HexToBytes failed, expected %s, actual %s\n", expected, hash)
	}
}

func TestFormatUrl(t *testing.T) {
	host := "localhost"
	var port uint64 = 8080
	expected := "localhost:8080"
	actual := FormatUrl(host, port)
	if expected != actual {
		t.Errorf("FormatWsUrl failed, expected %s, actual %s\n", expected, actual)
	}
}

func TestGenKeypair(t *testing.T) {
	secret := "test secret"

	hash := Hash(secret)
	expectedPrivateKey := ed25519.NewKeyFromSeed(hash)
	expectedPublicKey := expectedPrivateKey.Public().(PublicKey)

	actualPrivateKey, actualPublicKey := GenKeypair(secret)
	if !bytes.Equal(expectedPrivateKey, actualPrivateKey) {
		t.Errorf("GenKeypair failed, expected %s, actual %s\n", expectedPrivateKey, actualPrivateKey)
	}
	if !bytes.Equal(expectedPublicKey, actualPublicKey) {
		t.Errorf("GenKeypair failed, expected %s, actual %s\n", expectedPublicKey, actualPublicKey)
	}
}

func TestId(t *testing.T) {
	uuid1 := Id()
	uuid2 := Id()
	if uuid1 == uuid2 {
		t.Errorf("Id() failed, uuid1 == uuid2\n")
	}
}

func TestSign(t *testing.T) {
	secret := "test secret"
	data := "test data"
	hash := Hash(data)
	privateKey, _ := GenKeypair(secret)
	expected := ed25519.Sign(privateKey, hash)
	actual := Sign(privateKey, hash)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Sign failed, expected %s, actual %s\n", expected, actual)
	}
}

func TestVerify(t *testing.T) {
	secret := "test secret"
	data := "test data"
	hash := Hash(data)
	privateKey, publicKey := GenKeypair(secret)
	signature := Sign(privateKey, hash)
	if !Verify(publicKey, hash, signature) {
		t.Errorf("Verify failed, expected true, actual false\n")
	}
}
