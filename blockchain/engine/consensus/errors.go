package consensus

import "errors"

var (
	errInvalidPoW    = errors.New("proof of work does not match difficulty")
	errIncorrectPoW  = errors.New("block hash does not match calculated hash")
	errInvalidTxHash = errors.New("block transaction hash does not match calculated transaction hash")
)
