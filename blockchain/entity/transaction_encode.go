package entity

type TransactionEncoder interface {
	EncodeTransaction(transaction *Transaction) ([]byte, error)
	DecodeTransaction(input []byte) (*Transaction, error)

	EncodeOutput(output *Output) ([]byte, error)
	DecodeOutput(input []byte) (*Output, error)

	EncodeSignedInput(input *SignedInput) ([]byte, error)
	DecodeSignedInput(input []byte) (*SignedInput, error)

	EncodeUnsignedInput(input *UnsignedInput) ([]byte, error)
	DecodeUnsignedInput(input []byte) (*UnsignedInput, error)
}
