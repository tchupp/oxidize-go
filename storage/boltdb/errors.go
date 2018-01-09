package boltdb

const (
	BucketNotFoundError    = Error("bucket not found")
	BlockDataEmptyError    = Error("block data is empty")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
