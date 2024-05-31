package pbft2

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
	expected := string(hash[:])
	actual := Hash(data)

	if expected != actual {
		t.Errorf("Hash failed, expected %s, actual %s\n", expected, actual)
	}
}

func TestGenKeypair(t *testing.T) {
	secret := "test secret"

	hash := Hash(secret)
	seed := []byte(hash)
	expectedPrivateKey := ed25519.NewKeyFromSeed(seed)
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

func TestKey2Str(t *testing.T) {
	secret := "test secret"
	privateKey, publicKey := GenKeypair(secret)
	privateKeyStr := Key2Str(privateKey)
	publicKeyStr := Key2Str(publicKey)
	expectedPrivateKeyStr := hex.EncodeToString(privateKey)
	expectedPublicKeyStr := hex.EncodeToString(publicKey)
	if expectedPrivateKeyStr != privateKeyStr {
		t.Errorf("Key2Str failed, expected %s, actual %s\n", expectedPrivateKeyStr, privateKeyStr)
	}
	if expectedPublicKeyStr != publicKeyStr {
		t.Errorf("Key2Str failed, expected %s, actual %s\n", expectedPublicKeyStr, publicKeyStr)
	}
}

func TestStr2Key(t *testing.T) {
	secret := "test secret"
	privateKey, publicKey := GenKeypair(secret)
	privateKeyStr := Key2Str(privateKey)
	publicKeyStr := Key2Str(publicKey)
	actualPrivateKey := Str2Key(privateKeyStr)
	actualPublicKey := Str2Key(publicKeyStr)
	if !bytes.Equal(privateKey, actualPrivateKey) {
		t.Errorf("Str2Key failed, expected %s, actual %s\n", privateKey, actualPrivateKey)
	}
	if !bytes.Equal(publicKey, actualPublicKey) {
		t.Errorf("Str2Key failed, expected %s, actual %s\n", publicKey, actualPublicKey)
	}
}

func TestSign(t *testing.T) {
	secret := "test secret"
	data := "test data"
	hash := Hash(data)
	privateKey, _ := GenKeypair(secret)
	expected := string(ed25519.Sign(privateKey, []byte(hash)))
	actual := Sign(privateKey, hash)

	if expected != actual {
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
