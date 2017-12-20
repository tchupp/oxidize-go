package blockchain

type Repository interface {
	Head() (*Block, error)

	Block([]byte) (*Block, error)

	SaveBlock(*Block) error
}
