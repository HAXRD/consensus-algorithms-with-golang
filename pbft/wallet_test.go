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
