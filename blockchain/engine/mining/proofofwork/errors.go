package proofofwork

const (
	MaxNonceOverflowError  = Error("max nonce hit with no solution")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
