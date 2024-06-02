package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

/**
A Block stores the pool collected from tx pool, featured with the following methods:
1. PrintBlock
2. Genesis
3. HashBlock
4. VerifyBlock
5. VerifyBlockProposer
*/

type Block struct {
	Timestamp      string        `json:"timestamp"`
	LastHash       string        `json:"lastHash"`
	Data           []Transaction `json:"data"`
	Hash           string        `json:"hash"`
	Proposer       PublicKey     `json:"proposer"`
	Signature      string        `json:"signature"`
	Nonce          uint64        `json:"nonce"`
	PrePrepareMsgs []Message     `json:"prePrepareMsgs"`
	PrepareMsgs    []Message     `json:"prepareMsgs"`
	CommitMsgs     []Message     `json:"commitMsgs"`
}

/**
BlockPool stores pool temporarily proposed by proposer for each node.
It features the following methods:
1. NewBlockPool
2. BlockExists
3. AddBlock2Pool
4. GetBlock
5. CleanBlock
*/

type BlockPool struct {
	pool []Block
}

// PrintBlock prints the block info
func PrintBlock(block Block) {
	fmt.Printf("Block - "+
		"Timestamp: %s\n"+
		"LastHash: %s\n"+
		"Data: %v\n"+
		"Hash: %s\n"+
		"Nonce: %d\n"+
		"Signature: %s\n",
		block.Timestamp,
		block.LastHash,
		block.Data,
		block.Hash,
		block.Nonce,
		block.Signature,
	)
}

// Genesis creates the first block of the chain
func Genesis() *Block {
	return &Block{
		Timestamp:      time.Now().String(),
		LastHash:       "-",
		Data:           nil,
		Hash:           "-",
		Proposer:       nil,
		Signature:      "-",
		Nonce:          0,
		PrePrepareMsgs: nil,
		PrepareMsgs:    nil,
		CommitMsgs:     nil,
	}
}

// HashBlock hashes the block with timestamp, lastBlock's hash, marshalled data and current nonce
func HashBlock(timestamp string, lastHash string, data []Transaction, nonce uint64) string {
	dataInByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return chain_util2.Hash(timestamp + lastHash + string(dataInByte) + strconv.FormatUint(nonce, 10))
}

// VerifyBlock verifies the block information and its signature
func VerifyBlock(block Block) bool {
	return chain_util2.Verify(
		block.Proposer,
		HashBlock(block.Timestamp, block.LastHash, block.Data, block.Nonce),
		block.Signature,
	)
}

// VerifyBlockProposer verifies the block's proposer matches the proposer
// TODO: might remove this
func VerifyBlockProposer(block Block, proposer PublicKey) bool {
	return chain_util2.Key2Str(block.Proposer) == chain_util2.Key2Str(proposer)
}

// NewBlockPool creates a new block pool
func NewBlockPool() *BlockPool {
	return &BlockPool{pool: make([]Block, 0)}
}

// BlockExists checks if a given block exists in the pool
// by comparing its hash
func (bp *BlockPool) BlockExists(hash string) (bool, int) {
	for idx, b := range bp.pool {
		if b.Hash == hash {
			return true, idx
		}
	}
	return false, -1
}

// AddBlock2Pool adds a block to the block pool
func (bp *BlockPool) AddBlock2Pool(block Block) {
	bp.pool = append(bp.pool, block)
	log.Printf("Added block [%s] to pool\n", chain_util2.FormatHash(block.Hash))
}

// GetBlock get a copy of the block from the pool with given hash
func (bp *BlockPool) GetBlock(hash string) *Block {
	for _, b := range bp.pool {
		if b.Hash == hash {
			blockCopy := b
			return &blockCopy
		}
	}
	return nil
}

// CleanPool removes the block from the pool by matching block hash
func (bp *BlockPool) CleanPool(hash string) bool {
	exists, idx := bp.BlockExists(hash)
	if exists {
		bp.pool = append(bp.pool[:idx], bp.pool[idx+1:]...)
		return true
	} else {
		return false
	}
}
