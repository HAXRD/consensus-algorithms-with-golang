package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
)

type Blockchain struct {
	validators    []PublicKey
	chain         []Block
	blockIndexMap map[string]int
}

func NewBlockchain(vs Validators) *Blockchain {
	validators := vs.list
	chain := make([]Block, 0, 1)
	chain = append(chain, *Genesis())
	blockIndexMap := make(map[string]int)
	blockIndexMap[pow_util.Byte2Hex(chain[0].Hash)] = 0
	return &Blockchain{
		validators:    validators,
		chain:         chain,
		blockIndexMap: blockIndexMap,
	}
}

func (bc *Blockchain) VerifyBlock(block Block) bool {
	if VerifyBlock(block) {
		return true
	} else {
		return false
	}
}

func (bc *Blockchain) VerifyBlockWithLastBlockInChain(block Block) bool {
	lastBlock := bc.chain[len(bc.chain)-1]
	return pow_util.Byte2Hex(block.LastHash) == pow_util.Byte2Hex(lastBlock.Hash)
}

func (bc *Blockchain) AddBlock(block Block) {
	bc.chain = append(bc.chain, block)
	bc.blockIndexMap[pow_util.Byte2Hex(block.Hash)] = len(bc.chain) - 1
}

func (bc *Blockchain) BlockConflicts(block Block) bool {
	idx, ok := bc.blockIndexMap[pow_util.Byte2Hex(block.LastHash)]
	return ok && idx != len(bc.chain)-1
}

func (bc *Blockchain) Overwritable(block Block) bool {
	index, ok := bc.blockIndexMap[pow_util.Byte2Hex(block.LastHash)]
	if ok && index+1 <= len(bc.chain)-1 {
		targetBlock := bc.chain[index+1]
		return block.Timestamp < targetBlock.Timestamp
	}
	return false
}

func (bc *Blockchain) OverwriteBlock(block Block) []Block {
	index, ok := bc.blockIndexMap[pow_util.Byte2Hex(block.LastHash)]
	if ok && index+1 <= len(bc.chain)-1 {
		// delete pairs in map
		for i := index + 1; i < len(bc.chain); i++ {
			curBlock := bc.chain[i]
			delete(bc.blockIndexMap, pow_util.Byte2Hex(curBlock.Hash))
		}
		// store overwritten blocks
		blocksThatWereOverwritten := bc.chain[index+1:]
		// delete following overwritten blocks in chain
		bc.chain = bc.chain[:index+1]
		// add block to the chain
		bc.AddBlock(block)
		return blocksThatWereOverwritten
	}
	return nil
}

func (bc *Blockchain) BlockExists(block Block) bool {
	_, ok := bc.blockIndexMap[pow_util.Byte2Hex(block.Hash)]
	return ok
}

func (bc *Blockchain) Clear() {
	bc.chain = bc.chain[:1]
	bc.blockIndexMap = make(map[string]int)
	bc.blockIndexMap[pow_util.Byte2Hex(bc.chain[0].Hash)] = 0
}
