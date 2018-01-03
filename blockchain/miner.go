package blockchain

type MiningService interface {
	MineBlock(header *BlockHeader) *Block
}
