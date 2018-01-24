package node

import (
	"net"
	"testing"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

func TestBaseNode_AddPeer(t *testing.T) {
	bc, err := blockchain.Open(memdb.NewBlockRepository(), nil)
	if err != nil {
		t.Fatalf("opening blockchain: %s", err)
	}

	expectedHeader, err := bc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting best header with blockchain: %s", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	node1 := NewNode(bc, rpc.NewServer(lis))
	node1.Serve()
	defer node1.Shutdown()

	node2 := NewNode(bc, rpc.NewServer(nil))

	if len(node2.ActivePeers()) != 0 {
		t.Fatalf("incorrect starting peer count. got - %d, wanted  - %d", len(node2.ActivePeers()), 0)
	}

	if err := node2.AddPeer(lis.Addr().String()); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	activePeers := node2.ActivePeers()
	if len(activePeers) != 1 {
		t.Fatalf("incorrect final peer count. got - %d, wanted - %d", len(activePeers), 1)
	}

	peer := activePeers[0]
	if peer.Address != lis.Addr().String() {
		t.Errorf("incorrect peer address. got - %s, wanted - %s", peer.Address, lis.Addr())
	}
	if !peer.Head.IsEqual(expectedHeader.Hash) {
		t.Errorf("incorrect peer address. got - %s, wanted - %s", peer.Head, expectedHeader.Hash)
	}
}

func TestBaseNode_AddPeer_TargetIsOffline(t *testing.T) {
	bc, err := blockchain.Open(memdb.NewBlockRepository(), nil)
	if err != nil {
		t.Fatalf("opening blockchain: %s", err)
	}

	node := NewNode(bc, rpc.NewServer(nil))

	if len(node.ActivePeers()) != 0 {
		t.Fatalf("incorrect starting peer count. got - %d, wanted  - %d", len(node.ActivePeers()), 0)
	}

	if err := node.AddPeer("127.0.0.1:0"); err == nil {
		t.Fatal("expected error, got none")
	}

	activePeers := node.ActivePeers()
	if len(activePeers) != 0 {
		t.Fatalf("incorrect final peer count. got - %d, wanted - %d", len(activePeers), 0)
	}
}
