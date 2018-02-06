package node

import (
	"fmt"
	"time"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/blockrpc"
	"github.com/tclchiam/oxidize-go/p2p"
)

func startSyncHeadersFlow(peer *p2p.Peer, peerManager p2p.PeerManager, bc blockchain.Blockchain) {
	syncLogger := log.WithField("peer", peer.Address)

	start := time.Now()
	syncLogger.Info("starting sync")

	err := syncHeadersWithPeer(peer, peerManager, bc)

	syncLogger = syncLogger.WithField("elapsed", time.Since(start))
	if err != nil {
		syncLogger.WithError(err).Warn("error syncing with peer")
		return
	}

	syncLogger.Info("successfully synced with peer")
}

func syncHeadersWithPeer(peer *p2p.Peer, peerManager p2p.PeerManager, bc blockchain.Blockchain) error {
	conn := peerManager.GetPeerConnection(peer)
	if conn == nil {
		return fmt.Errorf("no connection open for peer")
	}

	syncClient := blockrpc.NewSyncClient(conn)

	for {
		peerBestHeader, err := syncClient.GetBestHeader()
		if err != nil {
			return fmt.Errorf("error getting peer best header: %s", err)
		}

		localBestHeader, err := bc.GetBestHeader()
		if err != nil {
			return fmt.Errorf("error getting local best header")
		}

		if localBestHeader.Index >= peerBestHeader.Index {
			return nil
		}

		peerHeaders, err := syncClient.GetHeaders(localBestHeader.Hash, localBestHeader.Index)
		if err != nil {
			return fmt.Errorf("error getting peer headers: %s", err)
		}

		if err = bc.SaveHeaders(peerHeaders); err != nil {
			return fmt.Errorf("error saving headers: %s", err)
		}

		if err := reconcileBlocks(syncClient, bc); err != nil {
			return err
		}
	}
}
