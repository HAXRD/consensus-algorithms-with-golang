package pbft

import (
	"encoding/hex"
	"fmt"
)

/**
BlockPool stores block temporarily. It holds the blocks until
it is added to the chain.
A block is added to the block pool when a `PRE-PREPARE` message
is received.

1. NewBlockPool: create a block pool
2. BlockExists: check if a given block exists or not
3. AddBlock2Pool: add a given block to the block pool
4. GetBlock: get block from the pool by hash
*/

type BlockPool struct {
	list []Block
}

func NewBlockPool() *BlockPool {
	return &BlockPool{make([]Block, 0)}
}

func (bp *BlockPool) BlockExists(block Block) bool {
	for _, b := range bp.list {
		if b.hash == block.hash {
			return true
		}
	}
	return false
}

func (bp *BlockPool) AddBlock2Pool(block Block) {
	bp.list = append(bp.list, block)
	fmt.Printf("Added block %s to pool\n", hex.EncodeToString([]byte(block.hash))[:5])
}

func (bp *BlockPool) GetBlock(hash string) *Block {
	for _, b := range bp.list {
		if b.hash == hash {
			return &b
		}
	}
	return nil
}
