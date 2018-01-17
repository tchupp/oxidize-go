package encoding

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/golang/protobuf/proto"
)

func toBlockData(block *entity.Block) *Block {
	var transactions []*Transaction
	for _, transaction := range block.Transactions() {
		transactions = append(transactions, toTransactionData(transaction))
	}

	return &Block{
		Header:       toBlockHeaderData(block.Header()),
		Transactions: transactions,
	}
}

func fromBlockData(block *Block) (*entity.Block, error) {
	var transactions []*entity.Transaction
	for _, txData := range block.Transactions {
		transaction, err := fromTransactionData(txData)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	header, err := fromBlockHeaderData(block.Header)
	if err != nil {
		return nil, err
	}

	return entity.NewBlock(
		header,
		transactions,
	), nil
}

func toBlockHeaderData(header *entity.BlockHeader) *BlockHeader {
	return &BlockHeader{
		Index:            proto.Uint64(header.Index),
		PreviousHash:     header.PreviousHash.Slice(),
		Timestamp:        proto.Uint64(header.Timestamp),
		TransactionsHash: header.TransactionsHash.Slice(),
		Nonce:            proto.Uint64(header.Nonce),
		Hash:             header.Hash.Slice(),
	}
}

func fromBlockHeaderData(header *BlockHeader) (*entity.BlockHeader, error) {
	previousHash, err := entity.NewHash(header.GetPreviousHash())
	if err != nil {
		return nil, fmt.Errorf("parsing previous hash: %s", err)
	}

	transactionsHash, err := entity.NewHash(header.GetTransactionsHash())
	if err != nil {
		return nil, fmt.Errorf("parsing transactions hash: %s", err)
	}

	hash, err := entity.NewHash(header.GetHash())
	if err != nil {
		return nil, fmt.Errorf("parsing block hash: %s", err)
	}

	return &entity.BlockHeader{
		Index:            header.GetIndex(),
		PreviousHash:     previousHash,
		Timestamp:        header.GetTimestamp(),
		TransactionsHash: transactionsHash,
		Nonce:            header.GetNonce(),
		Hash:             hash,
	}, nil
}
