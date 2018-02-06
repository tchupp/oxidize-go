package storage

import (
	"github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/logger"
)

var log = logger.Disabled

func UseLogger(logger *logrus.Entry) {
	log = logger
}

type loggingStorageDecorator struct {
	entity.ChainRepository
}

func WrapWithLogger(repository entity.ChainRepository) entity.ChainRepository {
	return &loggingStorageDecorator{ChainRepository: repository}
}

func (s loggingStorageDecorator) SaveBlock(block *entity.Block) error {
	blockLogger := log.WithField("block", block.Hash())

	err := s.ChainRepository.SaveBlock(block)
	if err != nil {
		blockLogger.WithError(err).Warn("failed to save block")
		return err
	}

	blockLogger.Debug("saved block")
	return nil
}

func (s loggingStorageDecorator) SaveHeader(header *entity.BlockHeader) error {
	headerLogger := log.WithField("header", header.Hash)

	err := s.ChainRepository.SaveHeader(header)
	if err != nil {
		headerLogger.WithError(err).Warn("failed to save header")
		return err
	}

	headerLogger.Debug("saved header")
	return nil
}
