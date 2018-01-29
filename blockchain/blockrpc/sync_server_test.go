package blockrpc

import (
	"fmt"
	"net"
	"math/rand"
	"testing"
	"google.golang.org/grpc"

	"github.com/google/go-cmp/cmp"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

func TestSyncServer_GetBestHeader(t *testing.T) {
	bc, err := blockchain.Open(memdb.NewBlockRepository(), memdb.NewHeaderRepository(), nil)
	if err != nil {
	}

	expectedHeader, err := bc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting best header: %s", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	server := rpc.NewServer(lis)
	server.RegisterSyncServer(NewSyncServer(bc))
	server.Serve()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewSyncClient(conn)

	header, err := client.GetBestHeader()
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	if !header.IsEqual(expectedHeader) {
		t.Errorf("headers don't match. got - %s, wanted - %s", header, expectedHeader)
	}
}

func TestSyncServer_GetHeaders(t *testing.T) {
	blockRepository := memdb.NewBlockRepository()
	bc, err := blockchain.Open(blockRepository, memdb.NewHeaderRepository(), nil)
	if err != nil {
		t.Fatalf("opening blockchain: %s", err)
	}

	saveRandomBlocks(bc, rand.Intn(31))

	expectedHeaders, err := bc.GetHeaders(nil, 0)
	if err != nil {
		t.Fatalf("getting headers: %s", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	server := rpc.NewServer(lis)
	server.RegisterSyncServer(NewSyncServer(bc))
	server.Serve()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewSyncClient(conn)

	bestHeader, err := bc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting best header: %s", err)
	}

	headers, err := client.GetHeaders(bestHeader.Hash, bestHeader.Index)
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	if !cmp.Equal(headers, expectedHeaders) {
		t.Errorf("headers don't match. got - %s, wanted - %s", headers, expectedHeaders)
	}
}

func saveRandomBlocks(bc blockchain.Blockchain, num int) error {
	miner := proofofwork.NewDefaultMiner(identity.RandomIdentity())

	for i := 0; i < num; i++ {
		coinbase := identity.RandomIdentity()
		head, err := bc.GetBestHeader()
		if err != nil {
			return fmt.Errorf("error reading best header")
		}

		transactions := entity.Transactions{entity.NewCoinbaseTx(coinbase, encoding.TransactionProtoEncoder())}
		block := miner.MineBlock(head, transactions)
		bc.SaveBlock(block)
	}

	return nil
}
