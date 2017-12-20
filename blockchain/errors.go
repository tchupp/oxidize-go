package blockchain

const (
	LatestHashNotFoundError = Error("latest hash not found")
	BucketNotFoundError     = Error("bucket not found")
	BlockDataEmptyError     = Error("block data is empty")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
