package node

import (
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage/memdb"
)

func TestBaseNode_AddPeer(t *testing.T) {
	remoteBc := buildBlockchain(t)
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis))
	remoteNode.Serve()
	defer remoteNode.Shutdown()

	localBc := buildBlockchain(t)
	localNode := newNode(localBc, rpc.NewServer(nil))

	verifyPeerCount(localNode, 0, t)

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	verifyPeerCount(localNode, 1, t)

	time.Sleep(600 * time.Millisecond)

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

	localBc := buildBlockchain(t)
	localNode := newNode(localBc, rpc.NewServer(nil))

	verifyPeerCount(localNode, 0, t)

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	verifyPeerCount(localNode, 1, t)

	remoteNode.Shutdown()

	time.Sleep(600 * time.Millisecond)

	verifyPeerCount(localNode, 0, t)
}

func TestBaseNode_AddPeer_TargetIsOffline(t *testing.T) {
	localBc := buildBlockchain(t)
	localNode := newNode(localBc, rpc.NewServer(nil))

	verifyPeerCount(localNode, 0, t)

	if _, err := localNode.AddPeer("127.0.0.1:0"); err == nil {
		t.Fatal("expected error, got none")
	}

	verifyPeerCount(localNode, 0, t)

	time.Sleep(600 * time.Millisecond)
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

	localBestBlock, err := localBc.GetBestBlock()
	if err != nil {
		t.Fatalf("getting local best block: %s", err)
	}
	remoteBestBlock, err := remoteBc.GetBestBlock()
	if err != nil {
		t.Fatalf("getting remote best block: %s", err)
	}
	if !remoteBestBlock.IsEqual(localBestBlock) {
		t.Errorf("unexpected local best block. got - %s, wanted - %s", localBestBlock, remoteBestBlock)
	}
}

func buildBlockchain(t *testing.T) blockchain.Blockchain {
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

func verifyPeerCount(node *baseNode, peerCount int, t *testing.T) {
	if len(node.ActivePeers()) != peerCount {
		t.Fatalf("incorrect peer count. got - %d, wanted  - %d", len(node.ActivePeers()), peerCount)
	}
}
