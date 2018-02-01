package boltdb

const (
	BucketNotFoundError    = Error("bucket not found")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
