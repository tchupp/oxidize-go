package node

import (
	"fmt"
	"time"

	"github.com/tclchiam/block_n_go/blockchain/blockrpc"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/p2p"
)

type syncBackend interface {
	GetBestHeader() (*entity.BlockHeader, error)
	SaveHeaders(headers entity.BlockHeaders) error
}

func startSyncFlow(peer *p2p.Peer, peerManager p2p.PeerManager, backend syncBackend) {
	syncLogger := log.WithField("peer", peer.Address)

	start := time.Now()
	syncLogger.Info("starting sync")

	err := syncWithPeer(peer, peerManager, backend)

	syncLogger = syncLogger.WithField("elapsed", time.Since(start))
	if err != nil {
		syncLogger.WithError(err).Warn("error syncing with peer")
		return
	}

	syncLogger.Info("successfully synced with peer")
	return
}

func syncWithPeer(peer *p2p.Peer, peerManager p2p.PeerManager, backend syncBackend) error {
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

		localBestHeader, err := backend.GetBestHeader()
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

		if len(peerHeaders) == 0 {
			return nil
		}

		if err = backend.SaveHeaders(peerHeaders); err != nil {
			return fmt.Errorf("error saving headers: %s", err)
		}
	}
}
