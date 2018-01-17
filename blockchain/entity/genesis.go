package entity

func DefaultGenesisBlock() *Block {
	return &Block{
		header: &BlockHeader{
			Index:            0,
			PreviousHash:     &EmptyHash,
			Timestamp:        uint64(1516146240),
			TransactionsHash: NewHashOrPanic("85fca3e06fe7196148c3f6beae1aeb7dc8b9308ce26bbd0f32bda91a60d63bbc"),
			Nonce:            80212,
			Hash:             NewHashOrPanic("0000129102d03577d79c9b919dcc71d99f7cf6c6b68f5e1f77d5adbfc0103938"),
		},
		transactions: Transactions{
			{
				ID: NewHashOrPanic("d81f935f0c45cd0df0ccf073ae0e33432dd14cd925262a51a5bb43a77f433862"),
				Outputs: []*Output{
					{Index: 0, Value: 10, PublicKeyHash: []byte("75a53812febc119e562265a250552182274fd970")},
				},
				Secret: []byte("39f39efae3884e28f5a5c4a62dd994e2943ac9cc7f9684070dfd4add7353722f"),
			},
		},
	}
}
