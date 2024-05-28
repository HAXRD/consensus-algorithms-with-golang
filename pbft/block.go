package pbft

import (
	"consensus-algorithms-with-golang/pbft/chainutil"
	"encoding/json"
	"fmt"
	"time"
)

/**
The block has the following properties:
1. timestamp: the time at which the block was created
2. lastHash: hash of the last block
3. hash: hash of the current block
4. data: the transaction data current block holds
5. proposer: the publicKey of the creator of the block
6. signature: the signed hash of the block
7. sequenceNo: the sequence num of the block
*/

type Block struct {
	timestamp       string        `json:"timestamp"`
	lastHash        string        `json:"lastHash"`
	hash            string        `json:"hash"`
	data            []Transaction `json:"data"`
	proposer        string        `json:"proposer"`
	signature       string        `json:"signature"`
	sequenceNum     uint64        `json:"sequenceNum"`
	prepareMessages []Prepare     `json:"prepareMessages"`
	commitMessages  []Commit      `json:"commitMessages"`
}

// NewBlock creates a new block with given data
func NewBlock(
	timestamp string,
	lastHash string,
	hash string,
	data []Transaction,
	proposer string,
	signature string,
	sequenceNum uint64) *Block {
	return &Block{
		timestamp,
		lastHash,
		hash,
		data,
		proposer,
		signature,
		sequenceNum,
		make([]Prepare, 0),
		make([]Commit, 0)}
}

// PrintBlock prints the given block
func PrintBlock(block Block) {
	fmt.Printf(
		"Block--->"+
			"timestamp:\t%s\n"+
			"lastHash:\t%s\n"+
			"hash:\t%s\n"+
			"data:\t%s\n"+
			"proposer:\t%s\n"+
			"signature:\t%s\n"+
			"sequenceNum:\t%d\n",
		block.timestamp,
		block.lastHash,
		block.hash,
		block.data,
		block.proposer,
		block.signature,
		block.sequenceNum)
}

// Genesis creates the very first block for the chain
func Genesis() Block {
	return Block{
		"genesis time",
		"----",
		"genesis-hash",
		[]Transaction{},
		"====",
		"SIGN",
		0,
		nil,
		nil}
}

// HashBlock hashes the given block info with timestamp, lastHash and data
func HashBlock(timestamp string, lastHash string, data []Transaction) string {
	dataInByteSlice, _ := json.Marshal(data)
	dataInString := string(dataInByteSlice)
	return chainutil.Hash(timestamp + lastHash + dataInString)
}

// SignBlockHash signs the block hash and returns the signature
func SignBlockHash(hash string, wallet Wallet) string {
	return wallet.Sign(hash)
}

// CreateBlock creates block with lastBlock, data and wallet
func CreateBlock(lastBlock Block, data []Transaction, wallet Wallet) *Block {
	timestamp := time.Now().String()
	hash := HashBlock(timestamp, lastBlock.hash, data)
	proposer := chainutil.Key2Str(wallet.GetPublicKey())
	signature := SignBlockHash(hash, wallet)
	return &Block{
		timestamp,
		lastBlock.hash,
		hash,
		data,
		proposer,
		signature,
		lastBlock.sequenceNum + 1,
		make([]Prepare, 0),
		make([]Commit, 0)}
}

// VerifyBlock verifies the block information, making sure the raw data matches
// the signature
func VerifyBlock(block Block) bool {
	return chainutil.Verify(
		block.proposer,
		HashBlock(block.timestamp, block.lastHash, block.data),
		block.signature)
}

// VerifyProposer verifies the proposer of the block matches the suggested proposer
func VerifyProposer(block Block, proposer string) bool {
	return block.proposer == proposer
}
