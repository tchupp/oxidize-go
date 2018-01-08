package blockchain

const (
	HeadBlockNotFoundError = Error("latest hash not found")
	BucketNotFoundError    = Error("bucket not found")
	BlockDataEmptyError    = Error("block data is empty")
	MaxNonceOverflowError  = Error("max nonce hit with no solution")

	TransactionInputHasBadSignatureMessage = "transaction '%s' has input '%d' with bad signature"
)

type Error string

func (e Error) Error() string {
	return string(e)
}
