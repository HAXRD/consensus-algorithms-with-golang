package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

/**
A Block stores the pool collected from tx pool, featured with the following methods:
1. NewBlock
2. Genesis
3. HashBlock
4. VerifyBlock
5. VerifyBlockProposer
*/

type Block struct {
	Timestamp   string        `json:"timestamp"`
	LastHash    []byte        `json:"lastHash"`
	Hash        []byte        `json:"hash"`
	Data        []Transaction `json:"data"`
	Proposer    PublicKey     `json:"proposer"`
	Signature   []byte        `json:"signature"`
	Nonce       uint64        `json:"nonce"`
	BlockMsgs   []Block       `json:"blockMsgs"`
	PrepareMsgs []Message     `json:"prepareMsgs"`
	CommitMsgs  []Message     `json:"commitMsgs"`
	RCMsgs      []Message     `json:"rcMsgs"`
	MsgType     string        `json:"msgType"`
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

// NewBlock creates a new block
func NewBlock(
	timestamp string,
	lastHash []byte,
	hash []byte,
	data []Transaction,
	proposer PublicKey,
	signature []byte,
	nonce uint64,
	blockMsgs []Block,
	prepareMsgs []Message,
	commitMsgs []Message,
	rcMsgs []Message,
) *Block {
	block := &Block{
		Timestamp:   timestamp,
		LastHash:    lastHash,
		Hash:        hash,
		Data:        data,
		Proposer:    proposer,
		Signature:   signature,
		Nonce:       nonce,
		BlockMsgs:   blockMsgs,
		PrepareMsgs: prepareMsgs,
		CommitMsgs:  commitMsgs,
		RCMsgs:      rcMsgs,
		MsgType:     MsgPrePrepare,
	}
	return block
}

// Genesis creates the first block of the chain
func Genesis() *Block {
	return NewBlock(
		time.Now().String(),
		[]byte("------"),
		[]byte("------"),
		nil,
		[]byte("------"),
		[]byte("------"),
		0,
		nil,
		nil,
		nil,
		nil,
	)
}

// HashBlock returns the hash of the block with timestamp, lastBlock's hash,
// marshalled data and current nonce
func HashBlock(timestamp string, lastHash []byte, data []Transaction, nonce uint64) []byte {
	dataInByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return chain_util.Hash(
		timestamp + chain_util.BytesToHex(lastHash) + string(dataInByte) + strconv.FormatUint(nonce, 10),
	)
}

// VerifyBlock verifies the block information and its signature
func VerifyBlock(block Block) bool {
	hash := HashBlock(block.Timestamp, block.LastHash, block.Data, block.Nonce)
	if chain_util.BytesToHex(hash) != chain_util.BytesToHex(block.Hash) {
		return false
	}
	return chain_util.Verify(
		block.Proposer,
		block.Hash,
		block.Signature,
	)
}

// VerifyBlockProposer verifies the block's proposer matches the proposer
func VerifyBlockProposer(block Block, proposer PublicKey) bool {
	return chain_util.BytesToHex(block.Proposer) == chain_util.BytesToHex(proposer)
}

// NewBlockPool creates a new block pool
func NewBlockPool() *BlockPool {
	return &BlockPool{pool: make([]Block, 0)}
}

// BlockExists checks if a given block exists in the pool
// by comparing its hash
func (bp *BlockPool) BlockExists(hash []byte) (bool, int) {
	for idx, b := range bp.pool {
		if chain_util.BytesToHex(b.Hash) == chain_util.BytesToHex(hash) {
			return true, idx
		}
	}
	return false, -1
}

// AddBlock2Pool adds a block to the block pool
func (bp *BlockPool) AddBlock2Pool(block Block) bool {
	// skip if exists
	if exists, _ := bp.BlockExists(block.Hash); exists {
		return false
	}

	bp.pool = append(bp.pool, block)
	log.Printf("Added block [%s] to pool\n", chain_util.BytesToHex(block.Hash)[:6])
	return true
}

// GetBlock get a copy of the block from the pool with given hash
func (bp *BlockPool) GetBlock(hash []byte) *Block {
	for _, b := range bp.pool {
		if chain_util.BytesToHex(b.Hash) == chain_util.BytesToHex(hash) {
			blockCopy := b
			return &blockCopy
		}
	}
	return nil
}

// CleanPool removes the block from the pool by matching block hash
func (bp *BlockPool) CleanPool(hash []byte) bool {
	exists, idx := bp.BlockExists(hash)
	if exists {
		bp.pool = append(bp.pool[:idx], bp.pool[idx+1:]...)
		return true
	} else {
		return false
	}
}

// Clear clears contents of block pool
func (bp *BlockPool) Clear() {
	bp.pool = bp.pool[:0]
}
