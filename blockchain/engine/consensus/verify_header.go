package consensus

import (
	"github.com/tclchiam/block_n_go/blockchain/engine/mining"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func VerifyHeader(header *entity.BlockHeader) error {
	// TODO verify parent block exists
	if !mining.HasDifficulty(header.Hash, header.Difficulty) {
		return errInvalidPoW
	}
	if calculatedHash := mining.CalculateHeaderHash(header); !header.Hash.IsEqual(calculatedHash) {
		return errIncorrectPoW
	}
	return nil
}
