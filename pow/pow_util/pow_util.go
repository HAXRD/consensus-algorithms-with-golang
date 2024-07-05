package pow_util

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type PublicKey = ed25519.PublicKey
type PrivateKey = ed25519.PrivateKey

func Byte2Hex(src []byte) string {
	return hex.EncodeToString(src)
}

func Hex2Byte(src string) ([]byte, error) {
	res, err := hex.DecodeString(src)
	return res, err
}

func Id() string {
	uuidV1, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil.String()
	}
	return uuidV1.String()
}

func Hash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func Sign(priKey PrivateKey, hash []byte) []byte {
	return ed25519.Sign(priKey, hash)
}

func Verify(pubKey PublicKey, hash []byte, signature []byte) bool {
	return ed25519.Verify(pubKey, hash, signature)
}

func GenKeyPair(secret string) (PrivateKey, PublicKey) {
	hash := Hash(secret)

	priKey := ed25519.NewKeyFromSeed(hash)
	pubKey := priKey.Public().(PublicKey)

	return priKey, pubKey
}

func FormatUrl(host string, port uint64) string {
	return host + ":" + strconv.FormatUint(port, 10)
}

func HashMeetsDifficulty(hash []byte, difficulty uint64) bool {
	hashStr := Byte2Hex(hash)
	if len(hashStr) >= int(difficulty) &&
		hashStr[:difficulty] == strings.Repeat("0", int(difficulty)) {
		return true
	}
	return false
}
