package blockrpc

import (
	"net"
	"testing"
	"google.golang.org/grpc"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

func TestSyncServer_GetBestHeader(t *testing.T) {
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
