package blockchain

type Blockchain struct {
	Blocks []*CommittedBlock
}

func (bc *Blockchain) AddBlock(data string) *Blockchain {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash, prevBlock.Index+1)
	newBlocks := append(bc.Blocks, newBlock)
	return &Blockchain{newBlocks}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*CommittedBlock{NewGenesisBlock()}}
}
