package entity

const genesisDifficulty = 4

func DefaultGenesisBlock() *Block {
	return &Block{
		header: NewBlockHeader(
			0,
			&EmptyHash,
			NewHashOrPanic("85fca3e06fe7196148c3f6beae1aeb7dc8b9308ce26bbd0f32bda91a60d63bbc"),
			1516146240,
			84167,
			NewHashOrPanic("00007bd9a44c3fea74388c483c3fc2fc8ac593c3da5566fbc1427cbf023e4ed9"),
			genesisDifficulty,
		),
		transactions: Transactions{
			{
				ID: NewHashOrPanic("d81f935f0c45cd0df0ccf073ae0e33432dd14cd925262a51a5bb43a77f433862"),
				Outputs: []*Output{
					{Index: 0, Value: 10, PublicKeyHash: []byte{117, 165, 56, 18, 254, 188, 17, 158, 86, 34, 101, 162, 80, 85, 33, 130, 39, 79, 217, 112}},
				},
				Secret: []byte("39f39efae3884e28f5a5c4a62dd994e2943ac9cc7f9684070dfd4add7353722f"),
			},
		},
	}
}
