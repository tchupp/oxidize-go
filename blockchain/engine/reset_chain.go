package engine

import "github.com/tclchiam/oxidize-go/blockchain/entity"

func ResetGenesis(repository entity.ChainRepository) error {
	exists, err := genesisBlockExists(repository)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	genesisBlock := entity.DefaultGenesisBlock()
	return repository.SaveBlock(genesisBlock)
}

func genesisBlockExists(repository entity.ChainRepository) (bool, error) {
	head, err := repository.BestBlock()
	if err != nil {
		return false, err
	}
	if head == nil {
		return false, nil
	}
	return true, nil
}
