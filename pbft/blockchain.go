package pbft

import "fmt"

/**
Blockchain keeps a copy of distributed ledger for each node.
Note that, for pbft, blocks at same indices of the chain may have different
"prepare"&"commit" messages. However, blockHashes are the same since they
depend on the data contained in the block.

1. NewBlockchain
2. addBlock2Chain
3. CreateBlock
4. GetProposer
5. IsBlockValid
6. AddUpdatedBlock2Chain
*/

type Blockchain struct {
	validatorList []string
	chain         []Block
}

func NewBlockchain() *Blockchain {
	validatorList := NewValidators(NUM_OF_NODES).list
	chain := make([]Block, 0)
	chain = append(chain, Genesis())
	return &Blockchain{validatorList, chain}
}

func (bc *Blockchain) addBlock(b Block) {
	bc.chain = append(bc.chain, b)
	fmt.Println("NEW BLOCK ADDED TO CHAIN")
}

func (bc *Blockchain) CreateBlock(txs []Transaction, wallet Wallet) *Block {
	return CreateBlock(bc.chain[len(bc.chain)-1], txs, wallet)
}

func (bc *Blockchain) GetProposer() string {
	index := bc.chain[len(bc.chain)-1].hash[0] % NUM_OF_NODES
	return bc.validatorList[index]
}

func (bc *Blockchain) IsBlockValid(block Block) bool {
	lastBlock := bc.chain[len(bc.chain)-1]
	if block.lastHash == lastBlock.hash &&
		block.hash == HashBlock(block.timestamp, block.lastHash, block.data) &&
		VerifyBlock(block) &&
		VerifyProposer(block, block.proposer) {
		fmt.Println("Block is VALID")
		return true
	} else {
		fmt.Println("Block is NOT VALID")
		return false
	}
}

func (bc *Blockchain) AddUpdatedBlock2Chain(hash string, blockPool BlockPool, preparePool PreparePool, commitPool CommitPool) {
	block := blockPool.GetBlock(hash)
	block.prepareMessages = preparePool.mapOfList[hash]
	block.commitMessages = commitPool.mapOfList[hash]
	bc.addBlock(*block)
}
