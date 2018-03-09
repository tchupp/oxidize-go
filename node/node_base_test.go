package node

import (
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/testdata"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/server/httpserver"
	"github.com/tclchiam/oxidize-go/server/rpc"
)

func TestBaseNode_AddPeer(t *testing.T) {
	remoteBc := testdata.NewBlockchainBuilder(t).Build()
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis), httpserver.NewServer(""))
	remoteNode.Serve()
	defer remoteNode.Close()

	localBc := testdata.NewBlockchainBuilder(t).Build()
	localNode := newNode(localBc, rpc.NewServer(nil), httpserver.NewServer(""))

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

	expectedHeader, err := remoteBc.BestHeader()
	if err != nil {
		t.Fatalf("getting best header: %s", err)
	}
	if !peer.BestHash.IsEqual(expectedHeader.Hash) {
		t.Errorf("unexpected peer best header. got - %s, wanted - %s", peer.BestHash, expectedHeader.Hash)
	}
}

func TestBaseNode_AddPeer_PeerLoosesConnection(t *testing.T) {
	remoteBc := testdata.NewBlockchainBuilder(t).Build()
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis), httpserver.NewServer(""))
	remoteNode.Serve()

	localBc := testdata.NewBlockchainBuilder(t).Build()
	localNode := newNode(localBc, rpc.NewServer(nil), httpserver.NewServer(""))

	verifyPeerCount(localNode, 0, t)

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	verifyPeerCount(localNode, 1, t)

	remoteNode.Close()

	time.Sleep(600 * time.Millisecond)

	verifyPeerCount(localNode, 0, t)
}

func TestBaseNode_AddPeer_TargetIsOffline(t *testing.T) {
	localBc := testdata.NewBlockchainBuilder(t).Build()
	localNode := newNode(localBc, rpc.NewServer(nil), httpserver.NewServer(""))

	verifyPeerCount(localNode, 0, t)

	if _, err := localNode.AddPeer("127.0.0.1:0"); err == nil {
		t.Fatal("expected error, got none")
	}

	verifyPeerCount(localNode, 0, t)

	time.Sleep(600 * time.Millisecond)
}

func TestBaseNode_AddPeer_SyncsHeadersWithNewPeer_WhenPeersVersionIsHigher(t *testing.T) {
	remoteBc := testdata.NewBlockchainBuilder(t).Build()
	saveRandomBlocks(t, remoteBc, rand.Intn(12))
	lis := buildListener(t)

	remoteNode := newNode(remoteBc, rpc.NewServer(lis), httpserver.NewServer(""))
	remoteNode.Serve()
	defer remoteNode.Close()

	localBc := testdata.NewBlockchainBuilder(t).Build()
	localNode := newNode(localBc, rpc.NewServer(nil), httpserver.NewServer(""))

	if _, err := localNode.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	time.Sleep(500 * time.Millisecond)

	localBestHeader, err := localBc.BestHeader()
	if err != nil {
		t.Fatalf("getting local best header: %s", err)
	}
	remoteBestHeader, err := remoteBc.BestHeader()
	if err != nil {
		t.Fatalf("getting remote best header: %s", err)
	}
	if !remoteBestHeader.IsEqual(localBestHeader) {
		t.Errorf("unexpected local best header. got - %s, wanted - %s", localBestHeader, remoteBestHeader)
	}

	localBestBlock, err := localBc.BestBlock()
	if err != nil {
		t.Fatalf("getting local best block: %s", err)
	}
	remoteBestBlock, err := remoteBc.BestBlock()
	if err != nil {
		t.Fatalf("getting remote best block: %s", err)
	}
	if !remoteBestBlock.IsEqual(localBestBlock) {
		t.Errorf("unexpected local best block. got - %s, wanted - %s", localBestBlock, remoteBestBlock)
	}
}

func buildListener(t *testing.T) net.Listener {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}
	return lis
}

func saveRandomBlocks(t *testing.T, bc blockchain.Blockchain, num int) {
	miner := proofofwork.NewDefaultMiner(identity.RandomIdentity().Address())

	for i := 0; i < num; i++ {
		beneficiary := identity.RandomIdentity().Address()
		head, err := bc.BestHeader()
		if err != nil {
			t.Fatal("error reading best header")
		}

		transactions := entity.Transactions{entity.NewRewardTx(beneficiary, encoding.TransactionProtoEncoder())}
		block := miner.MineBlock(head, transactions)
		bc.SaveBlock(block)
	}
}

func verifyPeerCount(node *baseNode, peerCount int, t *testing.T) {
	if len(node.ActivePeers()) != peerCount {
		t.Fatalf("incorrect peer count. got - %d, wanted  - %d", len(node.ActivePeers()), peerCount)
	}
}
