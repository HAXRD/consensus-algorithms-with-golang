package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
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
2. CreateBlock
3. AddUpdatedBlock2Chain
4. GetProposer
5. VerifyBlock
*/

type Blockchain struct {
	validators []PublicKey
	chain      []Block
}

// NewBlockchain creates a new blockchain
func NewBlockchain(vs Validators) *Blockchain {
	validators := vs.list
	chain := make([]Block, 0, 1)
	chain = append(chain, *Genesis())
	return &Blockchain{
		validators: validators,
		chain:      chain,
	}
}

// CreateBlock creates a new block with given wallet and collected
// txs. It calls wallet's `CreateBlock` method.
func (bc *Blockchain) CreateBlock(wallet Wallet, txs []Transaction) *Block {
	return wallet.CreateBlock(bc.chain[len(bc.chain)-1], txs)
}

// AddUpdatedBlock2Chain first get a copy of block with given hash,
// then add the PRE-PREPARE pool, PREPARE pool, COMMIT pool to update
// the block. Finally, it adds the updated block to the chain.
func (bc *Blockchain) AddUpdatedBlock2Chain(
	hash []byte,
	blockPool BlockPool, preparePool MsgPool, commitPool MsgPool) {

	existsInPool, _ := blockPool.BlockExists(hash)
	if !existsInPool {
		log.Printf("Added block [%s] to blockchain failed, BLOCK NOT EXISTS IN BLOCK POOL!", chain_util.BytesToHex(hash)[:6])
	} else {
		hashHex := chain_util.BytesToHex(hash)
		block := blockPool.GetBlock(hash)
		// check if the proposed block is matching the lastblock
		if chain_util.BytesToHex(block.LastHash) != chain_util.BytesToHex(bc.chain[len(bc.chain)-1].Hash) {
			log.Printf("Added block [%s] to blockchain failed, BLOCK'S LASTHASH NOT MATCHED!", chain_util.BytesToHex(hash)[:6])
			return
		}

		block.BlockMsgs = blockPool.pool
		block.PrepareMsgs = preparePool.mapPool[hashHex]
		block.CommitMsgs = commitPool.mapPool[hashHex]
		bc.chain = append(bc.chain, *block)
		log.Printf("Added block [%s] to blockchain succeed!", chain_util.BytesToHex(hash)[:6])
	}
}

// GetProposer get the proposer according to the latest block's info in the chain
func (bc *Blockchain) GetProposer() PublicKey {
	index := bc.chain[len(bc.chain)-1].Hash[0] % NUM_OF_NODES
	return bc.validators[index]
}

// VerifyBlock verifies a block with respect to the blockchain
func (bc *Blockchain) VerifyBlock(block Block) bool {
	lastBlock := bc.chain[len(bc.chain)-1]
	if chain_util.BytesToHex(block.LastHash) == chain_util.BytesToHex(lastBlock.Hash) &&
		chain_util.BytesToHex(block.Hash) == chain_util.BytesToHex(HashBlock(block.Timestamp, block.LastHash, block.Data, block.Nonce)) &&
		VerifyBlock(block) &&
		VerifyBlockProposer(block, block.Proposer) {
		log.Printf("Block [%s] is VALID", chain_util.BytesToHex(block.Hash)[:6])
		return true
	} else {
		log.Printf("Block [%s] is INVALID", chain_util.BytesToHex(block.Hash)[:6])
		return false
	}
}
