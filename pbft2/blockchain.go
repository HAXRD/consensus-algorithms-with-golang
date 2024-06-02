package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"log"
)

/**
Blockchain keeps a copy of the distributed ledger at each node.
Node that, for PBFT, blocks at the same indices may have different
phase transition messages (i.e., "PRE-PREPARE", "PREPARE", "COMMIT").
However, block hashes are the same since they depend on the data
contained in the block.
Blockchain features the following methods:
1. NewBlockchain
2. AddUpdatedBlock2Chain
3. GetProposer
4. VerifyBlock
*/

type Blockchain struct {
	validators []string
	chain      []Block
}

// NewBlockchain creates a new blockchain
func NewBlockchain() *Blockchain {
	validators := NewValidators(NUM_OF_NODES).list
	chain := make([]Block, 0, 1)
	chain = append(chain, *Genesis())
	return &Blockchain{
		validators: validators,
		chain:      chain,
	}
}

// AddUpdatedBlock2Chain first get a copy of block with given hash,
// then add the PRE-PREPARE pool, PREPARE pool, COMMIT pool to update
// the block. Finally, it adds the updated block to the chain.
func (bc *Blockchain) AddUpdatedBlock2Chain(
	hash string,
	blockPool BlockPool,
	prePreparePool MsgPool, preparePool MsgPool, commitPool MsgPool) bool {

	exists, _ := blockPool.BlockExists(hash)
	if !exists {
		return false
	} else {
		block := blockPool.GetBlock(hash)
		block.PrePrepareMsgs = prePreparePool.mapPool[hash]
		block.PrepareMsgs = preparePool.mapPool[hash]
		block.CommitMsgs = commitPool.mapPool[hash]
		bc.chain = append(bc.chain, *block)
		return true
	}
}

// GetProposer get the proposer according to the latest block's info in the chain
func (bc *Blockchain) GetProposer() string {
	index := bc.chain[len(bc.chain)-1].Hash[0] % NUM_OF_NODES
	return bc.validators[index]
}

// VerifyBlock verifies a block with respect to the blockchain
func (bc *Blockchain) VerifyBlock(block Block) bool {
	lastBlock := bc.chain[len(bc.chain)-1]
	if block.LastHash == lastBlock.Hash &&
		block.Hash == HashBlock(block.Timestamp, block.LastHash, block.Data, block.Nonce) &&
		VerifyBlock(block) &&
		VerifyBlockProposer(block, block.Proposer) {
		log.Printf("Block [%s] is VALID", chain_util2.FormatHash(block.Hash))
		return true
	} else {
		log.Printf("Block [%s] is INVALID", chain_util2.FormatHash(block.Hash))
		return false
	}
}
