package blockchain

type Blockchain struct {
	latestHash []byte
}

func (bc *Blockchain) AddBlock(data string) *Blockchain {
	return &Blockchain{latestHash: []byte{}}
}

func (bc *Blockchain) LatestHash() []byte {
	return bc.latestHash
}

func New(latestHash []byte) *Blockchain {
	return &Blockchain{latestHash: latestHash}
}

