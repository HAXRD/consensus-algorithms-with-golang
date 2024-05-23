package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"testing"
)

func TestNewWallet(t *testing.T) {
	secret := "secret"
	wallet := NewWallet(secret)
	privateKey, publicKey := chainutil.GenKeypair(secret)
	if chainutil.Key2Str(privateKey) != chainutil.Key2Str(wallet.privateKey) {
		t.Errorf(
			"Private key not matched, want, %s, got, %s",
			chainutil.Key2Str(privateKey),
			chainutil.Key2Str(wallet.privateKey))
	}
	if chainutil.Key2Str(publicKey) != chainutil.Key2Str(wallet.publicKey) {
		t.Errorf(
			"Public key not matched, want, %s, got, %s",
			chainutil.Key2Str(privateKey),
			chainutil.Key2Str(wallet.publicKey))
	}
}

func TestSign(t *testing.T) {
	wallet := NewWallet("secret")
	data := "test data"
	hash := chainutil.Hash(data)
	expectedSignature := chainutil.Sign(chainutil.Key2Str(wallet.privateKey), hash)
	actualSignature := wallet.Sign(hash)
	if expectedSignature != actualSignature {
		t.Errorf("Signature not matched, want: %s, got: %s", expectedSignature, actualSignature)
	}
}

func TestVerify(t *testing.T) {
	wallet := NewWallet("secret")
	data := "test data"
	hash := chainutil.Hash(data)
	signature := chainutil.Sign(chainutil.Key2Str(wallet.privateKey), hash)
	if !wallet.Verify(hash, signature) {
		t.Errorf("wallet.Verify failed")
	}
}
