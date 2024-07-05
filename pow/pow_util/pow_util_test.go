package pow_util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestByte2Hex(t *testing.T) {
	data := []byte("test data")
	expected := hex.EncodeToString(data)
	actual := Byte2Hex(data)
	if expected != actual {
		t.Errorf("Byte2Hex expect: %v, actual: %v\n", expected, actual)
	}
}

func TestHex2Byte(t *testing.T) {
	data := []byte("test data")
	dataHex := Byte2Hex(data)
	actual, err := Hex2Byte(dataHex)
	if err != nil {
		t.Error(err)
	}
	if dataHex != Byte2Hex(actual) {
		t.Errorf("Hex2Byte expect: %v, actual: %v\n", data, actual)
	}
}

func TestHash(t *testing.T) {
	data := "test data"
	actual := Hash(data)
	mock := sha256.Sum256([]byte(data))
	expected := mock[:]
	if !bytes.Equal(actual, expected) {
		t.Errorf("Hash expect: %v, actual: %v\n", expected, actual)
	}
}

func Test_GenKeyPair_Sign_Verify(t *testing.T) {
	secret := "secret"
	data := "test data"
	hash := Hash(data)
	priKey, pubKey := GenKeyPair(secret)
	signature := Sign(priKey, hash)
	if !Verify(pubKey, hash, signature) {
		t.Error("GenKeyPair, Sign, Verify failed.")
	}
}
