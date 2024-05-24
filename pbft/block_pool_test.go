package pbft

import "testing"

func TestNewBlockPool(t *testing.T) {
	blockPool := NewBlockPool()
	if len(blockPool.list) != 0 {
		t.Errorf("NewBlockPool fail")
	}
	genesis := Genesis()
	blockPool.AddBlock2Pool(genesis)
	if len(blockPool.list) != 1 {
		t.Errorf("NewBlockPool fail")
	}
}

func TestBlockPool_BlockExists(t *testing.T) {
	blockPool := NewBlockPool()
	genesis := Genesis()
	if blockPool.BlockExists(genesis) {
		t.Errorf("BlockExists fail")
	}
	blockPool.AddBlock2Pool(genesis)
	if !blockPool.BlockExists(genesis) {
		t.Errorf("BlockExists fail")
	}
}

func TestBlockPool_AddBlock2Pool(t *testing.T) {
	blockPool := NewBlockPool()
	genesis := Genesis()
	if len(blockPool.list) != 0 {
		t.Errorf("AddBlock2Pool fail")
	}
	blockPool.AddBlock2Pool(genesis)
	if len(blockPool.list) != 1 {
		t.Errorf("AddBlock2Pool fail")
	}
}

func TestBlockPool_GetBlock(t *testing.T) {
	blockPool := NewBlockPool()
	genesis := Genesis()
	block1 := blockPool.GetBlock(genesis.hash)
	if block1 != nil {
		t.Errorf("GetBlock fail")
	}
	blockPool.AddBlock2Pool(genesis)
	block2 := blockPool.GetBlock(genesis.hash)
	if genesis.hash != block2.hash {
		t.Errorf("GetBlock fail")
	}
}
