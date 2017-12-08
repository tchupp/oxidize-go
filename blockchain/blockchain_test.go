package blockchain

import "testing"

func TestNewBlockchain_StartWithGenesisBlock(t *testing.T) {
	blockchain := NewBlockchain()

	if len(blockchain.Blocks) != 1 {
		t.Errorf("Expected new Blockchain to have a single block, the Genesis Block")
	}

	genesisBlock := blockchain.Blocks[0]

	if genesisBlock.Index != 0 {
		t.Errorf("Expected Genesis Block to have Index of 0")
	}
	if len(genesisBlock.PreviousHash) != 0 {
		t.Errorf("Expected Genesis Block to have PreviousHash with len of 0")
	}
}

func TestBlockchain_AddBlock_DoesNotModifyExistingBlockchain(t *testing.T) {
	secondBlockData := "Second block yo"

	blockchain := NewBlockchain()
	newBlockchain := blockchain.AddBlock(secondBlockData)

	if len(blockchain.Blocks) != 1 {
		t.Errorf("Expected starting Blockchain to have block count %d, was %d", 1, len(blockchain.Blocks))
	}

	if len(newBlockchain.Blocks) != 2 {
		t.Errorf("Expected new Blockchain to have block count %d, was %d", 2, len(blockchain.Blocks))
	}
	if string(newBlockchain.Blocks[1].Data) != secondBlockData {
		t.Errorf("Expected new Block to have data %d, was %d", secondBlockData, newBlockchain.Blocks[1].Data)
	}
}
