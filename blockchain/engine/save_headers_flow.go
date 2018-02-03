package engine

import (
	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/block_n_go/blockchain/engine/consensus"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type headerChain interface {
	GetBestHeader() (*entity.BlockHeader, error)
	GetHeader(hash *entity.Hash) (*entity.BlockHeader, error)
	SaveHeader(*entity.BlockHeader) error
}

func SaveHeaders(headers entity.BlockHeaders, chain headerChain) error {
	for _, header := range headers.Sort() {
		_, err := saveHeader(header, chain)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveHeader(header *entity.BlockHeader, chain headerChain) (bool, error) {
	currentBestHeader, err := chain.GetBestHeader()
	if err != nil {
		return false, err
	}

	if err := consensus.VerifyHeader(header); err != nil {
		return false, err
	}

	switch {
	case currentBestHeader.Index >= header.Index:
		return false, nil

	case currentBestHeader.Index+1 == header.Index:
		return false, chain.SaveHeader(header)

	case currentBestHeader.Index+1 < header.Index:
		log.Warn("Future header, we shouldn't get this situation")
		return true, nil

	default:
		log.Warn("Somehow, a header case wasn't handled...")
		return false, nil
	}
}
