package node

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/blockrpc"
)

func reconcileBlocks(syncClient blockrpc.SyncClient, bc blockchain.Blockchain) error {
	for {
		bestHeader, err := bc.BestHeader()
		if err != nil {
			return fmt.Errorf("error reading best header: %s", err)
		}
		bestBlock, err := bc.BestBlock()
		if err != nil {
			return fmt.Errorf("error reading best block: %s", err)
		}

		if bestHeader.IsEqual(bestBlock.Header()) {
			return nil
		}

		block, err := syncClient.GetBlock(bestHeader.Hash, bestHeader.Index)
		if err := bc.SaveBlock(block); err != nil {
			return fmt.Errorf("error saving block: %s", err)
		}
	}
	return nil
}
