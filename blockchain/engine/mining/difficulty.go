package mining

import (
	"math"
	"math/big"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

var (
	maxUint256 = new(big.Int).Lsh(big.NewInt(1), 256)
)

func FindDifficulty(hash *entity.Hash) uint64 {
	hashInt := new(big.Int).SetBytes(hash.Slice())

	for difficulty := uint64(1); difficulty < math.MaxUint64; difficulty += 1 {
		target := new(big.Int).Rsh(maxUint256, uint(4*difficulty))

		if cmp := hashInt.Cmp(target); cmp >= 0 {
			return difficulty - 1
		}
	}
	return 0
}

func HasDifficulty(hash *entity.Hash, difficulty uint64) bool {
	if hash == nil {
		return false
	}

	hashInt := new(big.Int).SetBytes(hash.Slice())
	target := new(big.Int).Rsh(maxUint256, uint(4*difficulty))
	return hashInt.Cmp(target) == -1
}
