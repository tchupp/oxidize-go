package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"fmt"
)

func toBlockData(block *entity.Block) *Block {
	var transactions []*Transaction
	for _, transaction := range block.Transactions() {
		transactions = append(transactions, toTransactionData(transaction))
	}

	return &Block{
		Header:       toBlockHeaderData(block.Header()),
		Transactions: transactions,
		Hash:         block.Hash().Slice(),
		Nonce:        proto.Uint64(block.Nonce()),
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

	hash, err := chainhash.NewHash(block.GetHash())
	if err != nil {
		return nil, err
	}

	solution := &entity.BlockSolution{
		Hash:  hash,
		Nonce: block.GetNonce(),
	}

	return entity.NewBlock(
		header,
		solution,
		transactions,
	), nil
}

func toBlockHeaderData(header *entity.BlockHeader) *BlockHeader {
	return &BlockHeader{
		Index:            proto.Uint64(header.Index),
		PreviousHash:     header.PreviousHash.Slice(),
		Timestamp:        proto.Uint64(header.Timestamp),
		TransactionsHash: header.TransactionsHash.Slice(),
	}
}

func fromBlockHeaderData(block *BlockHeader) (*entity.BlockHeader, error) {
	previousHash, err := chainhash.NewHash(block.GetPreviousHash())
	if err != nil {
		return nil, fmt.Errorf("parsing previous hash: %s", err)
	}

	transactionsHash, err := chainhash.NewHash(block.GetTransactionsHash())
	if err != nil {
		return nil, fmt.Errorf("parsing transactions hash: %s", err)
	}

	return &entity.BlockHeader{
		Index:            block.GetIndex(),
		PreviousHash:     previousHash,
		Timestamp:        block.GetTimestamp(),
		TransactionsHash: transactionsHash,
	}, nil
}
