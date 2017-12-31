package blockchain

type Repository interface {
	Head() (head *Block, err error)

	Block([]byte) (*Block, error)

	SaveBlock(*Block) error
}
