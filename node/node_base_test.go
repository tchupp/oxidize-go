package node

import (
	"net"
	"math/rand"
	"testing"
	"time"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/storage/memdb"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
)

func TestBaseNode_AddPeer(t *testing.T) {
	remoteBc := buildBlockchain(t)
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis))
	remoteNode.Serve()
	defer remoteNode.Shutdown()

	localNode := newNode(nil, rpc.NewServer(nil))

	if len(localNode.ActivePeers()) != 0 {
		t.Fatalf("incorrect starting peer count. got - %d, wanted  - %d", len(localNode.ActivePeers()), 0)
	}

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(localNode.ActivePeers()) != 1 {
		t.Fatalf("incorrect final peer count. got - %d, wanted - %d", len(localNode.ActivePeers()), 1)
	}

	peer := localNode.ActivePeers()[0]
	if peer.Address != lis.Addr().String() {
		t.Errorf("incorrect peer address. got - %s, wanted - %s", peer.Address, lis.Addr())
	}

	expectedHeader, err := remoteBc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting best header: %s", err)
	}
	if !peer.BestHash.IsEqual(expectedHeader.Hash) {
		t.Errorf("unexpected peer best header. got - %s, wanted - %s", peer.BestHash, expectedHeader.Hash)
	}
}

func TestBaseNode_AddPeer_PeerLoosesConnection(t *testing.T) {
	remoteBc := buildBlockchain(t)
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis))
	remoteNode.Serve()

	localNode := newNode(nil, rpc.NewServer(nil))

	if len(localNode.ActivePeers()) != 0 {
		t.Fatalf("incorrect starting peer count. got - %d, wanted  - %d", len(localNode.ActivePeers()), 0)
	}

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(localNode.ActivePeers()) != 1 {
		t.Fatalf("incorrect intermediate peer count. got - %d, wanted - %d", len(localNode.ActivePeers()), 1)
	}

	remoteNode.Shutdown()

	time.Sleep(600 * time.Millisecond)

	if len(localNode.ActivePeers()) != 0 {
		t.Fatalf("incorrect final peer count. got - %d, wanted - %d", len(localNode.ActivePeers()), 0)
	}
}

func TestBaseNode_AddPeer_TargetIsOffline(t *testing.T) {
	localNode := newNode(nil, rpc.NewServer(nil))

	if len(localNode.ActivePeers()) != 0 {
		t.Fatalf("incorrect starting peer count. got - %d, wanted  - %d", len(localNode.ActivePeers()), 0)
	}

	if _, err := localNode.AddPeer("127.0.0.1:0"); err == nil {
		t.Fatal("expected error, got none")
	}

	activePeers := localNode.ActivePeers()
	if len(activePeers) != 0 {
		t.Fatalf("incorrect final peer count. got - %d, wanted - %d", len(activePeers), 0)
	}
}

func TestBaseNode_AddPeer_SyncsHeadersWithNewPeer_WhenPeersVersionIsHigher(t *testing.T) {
	remoteBc := buildBlockchain(t)
	saveRandomBlocks(t, remoteBc, rand.Intn(12))
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis))
	remoteNode.Serve()
	defer remoteNode.Shutdown()

	localBc := buildBlockchain(t)
	localNode := newNode(localBc, rpc.NewServer(nil))

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	time.Sleep(500 * time.Millisecond)

	localBestHeader, err := localBc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting local best header: %s", err)
	}
	remoteBestHeader, err := remoteBc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting remote best header: %s", err)
	}
	if !remoteBestHeader.IsEqual(localBestHeader) {
		t.Errorf("unexpected local best header. got - %s, wanted - %s", localBestHeader, remoteBestHeader)
	}
}

func buildBlockchain(t *testing.T) (blockchain.Blockchain) {
	bc, err := blockchain.Open(memdb.NewChainRepository(), nil)
	if err != nil {
		t.Fatalf("opening blockchain: %s", err)
	}
	return bc
}

func buildListener(t *testing.T) net.Listener {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}
	return lis
}

func saveRandomBlocks(t *testing.T, bc blockchain.Blockchain, num int) {
	miner := proofofwork.NewDefaultMiner(identity.RandomIdentity())

	for i := 0; i < num; i++ {
		coinbase := identity.RandomIdentity()
		head, err := bc.GetBestHeader()
		if err != nil {
			t.Fatal("error reading best header")
		}

		transactions := entity.Transactions{entity.NewCoinbaseTx(coinbase, encoding.TransactionProtoEncoder())}
		block := miner.MineBlock(head, transactions)
		bc.SaveBlock(block)
	}
}
