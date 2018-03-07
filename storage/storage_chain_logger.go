package storage

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type loggingChainStorageDecorator struct {
	entity.ChainRepository
}

func WrapChainWithLogger(repository entity.ChainRepository) entity.ChainRepository {
	return &loggingChainStorageDecorator{ChainRepository: repository}
}

func (s *loggingChainStorageDecorator) SaveBlock(block *entity.Block) error {
	blockLogger := log.WithField("block", block.Hash())

	if err := s.ChainRepository.SaveBlock(block); err != nil {
		blockLogger.WithError(err).Warn("failed to save block")
		return err
	}

	blockLogger.Debug("saved block")
	return nil
}

func (s *loggingChainStorageDecorator) SaveHeader(header *entity.BlockHeader) error {
	headerLogger := log.WithField("header", header.Hash)

	if err := s.ChainRepository.SaveHeader(header); err != nil {
		headerLogger.WithError(err).Warn("failed to save header")
		return err
	}

	headerLogger.Debug("saved header")
	return nil
}
