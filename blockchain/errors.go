package blockchain

const (
	HeadBlockNotFoundError = Error("latest hash not found")
	BucketNotFoundError    = Error("bucket not found")
	BlockDataEmptyError    = Error("block data is empty")
	MaxNonceOverflowError  = Error("max nonce hit with no solution")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
