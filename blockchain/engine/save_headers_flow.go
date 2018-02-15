package engine

import (
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/blockchain/engine/consensus"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

type headerChain interface {
	BestHeader() (*entity.BlockHeader, error)
	HeaderByHash(hash *entity.Hash) (*entity.BlockHeader, error)
	SaveHeader(*entity.BlockHeader) error
}

func SaveHeaders(headers entity.BlockHeaders, chain headerChain) error {
	var result *multierror.Error

	for _, header := range headers.Sort() {
		_, err := saveHeader(header, chain)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func saveHeader(header *entity.BlockHeader, chain headerChain) (bool, error) {
	currentBestHeader, err := chain.BestHeader()
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
		logrus.Warn("Future header, we shouldn't get this situation")
		return true, nil

	default:
		logrus.Warn("Somehow, a header case wasn't handled...")
		return false, nil
	}
}
