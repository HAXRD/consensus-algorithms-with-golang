package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Block use `timestamp + lastHash + Data + nonce` to get `Hash`
type Block struct {
	Timestamp string        `json:"timestamp"`
	LastHash  []byte        `json:"lastHash"`
	Hash      []byte        `json:"hash"`
	Data      []Transaction `json:"data"`
	Proposer  PublicKey     `json:"proposer"`
	Signature []byte        `json:"signature"`
	Nonce     uint64        `json:"nonce"`
}

func (block Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	return json.Marshal(&struct {
		LastHash  string `json:"lastHash"`
		Hash      string `json:"hash"`
		Proposer  string `json:"proposer"`
		Signature string `json:"signature"`
		Nonce     string `json:"nonce"`
		Alias
	}{
		LastHash:  pow_util.Byte2Hex(block.LastHash),
		Hash:      pow_util.Byte2Hex(block.Hash),
		Proposer:  pow_util.Byte2Hex(block.Proposer),
		Signature: pow_util.Byte2Hex(block.Signature),
		Nonce:     fmt.Sprintf("%d", block.Nonce),
		Alias:     Alias(block),
	})
}

func (block *Block) UnmarshalJSON(data []byte) error {
	type Alias Block
	aux := &struct {
		LastHash  string `json:"lastHash"`
		Hash      string `json:"hash"`
		Proposer  string `json:"proposer"`
		Signature string `json:"signature"`
		Nonce     string `json:"nonce"`
		*Alias
	}{
		Alias: (*Alias)(block),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	lastHash, err := pow_util.Hex2Byte(aux.LastHash)
	if err != nil {
		return err
	}
	hash, err := pow_util.Hex2Byte(aux.Hash)
	if err != nil {
		return err
	}
	proposer, err := pow_util.Hex2Byte(aux.Proposer)
	if err != nil {
		return err
	}
	signature, err := pow_util.Hex2Byte(aux.Signature)
	if err != nil {
		return err
	}
	nonce, err := strconv.ParseUint(aux.Nonce, 10, 64)
	if err != nil {
		return err
	}
	block.LastHash = lastHash
	block.Hash = hash
	block.Proposer = proposer
	block.Signature = signature
	block.Nonce = nonce
	return nil
}

func Genesis() *Block {
	return &Block{
		Timestamp: time.Now().String(),
		LastHash:  []byte("------"),
		Hash:      []byte("-----"),
		Data:      nil,
		Proposer:  []byte("------"),
		Signature: []byte("------"),
		Nonce:     0,
	}
}

func HashBlock(lastHash []byte, data []Transaction, nonce uint64) []byte {
	dataInByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return pow_util.Hash(
		pow_util.Byte2Hex(lastHash) + string(dataInByte) + strconv.FormatUint(nonce, 10),
	)
}

func VerifyBlock(block Block) bool {
	hash := HashBlock(block.LastHash, block.Data, block.Nonce)
	if pow_util.Byte2Hex(hash) != pow_util.Byte2Hex(block.Hash) {
		return false
	}
	return pow_util.Verify(
		block.Proposer,
		block.Hash,
		block.Signature,
	)
}
