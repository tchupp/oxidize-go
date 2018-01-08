package entity

type BlockEncoder interface {
	EncodeBlock(block *Block) ([]byte, error)
	DecodeBlock(input []byte) (*Block, error)
}
