package main

type Blockchain struct {
	blocks []*CommittedBlock
}

func (bc *Blockchain) AddBlock(data string) *Blockchain {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	newBlocks := append(bc.blocks, newBlock)
	return &Blockchain{newBlocks}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*CommittedBlock{NewGenesisBlock()}}
}
