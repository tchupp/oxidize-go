package p2p

import (
	"net"
	"strings"
	"testing"
	"time"
	"google.golang.org/grpc"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

func TestDiscoveryServer_Ping(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	server := rpc.NewServer(lis)
	server.RegisterDiscoveryServer(NewDiscoveryServer(nil))
	server.Serve()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTimeout(500*time.Millisecond), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewDiscoveryClient(conn)

	err = client.Ping()

	if err != nil {
		t.Errorf("expected no error, was: %s", err)
	}
}

func TestDiscoveryServer_Ping_TargetIsOffline(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:0", grpc.WithTimeout(500*time.Millisecond), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewDiscoveryClient(conn)

	err = client.Ping()

	expectedErrorMessage := "rpc error: code = Unavailable desc = grpc: the connection is unavailable"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("unexpected error message. got - %s, wanted: %s", err, expectedErrorMessage)
	}
}

func TestDiscoveryServer_Version(t *testing.T) {
	bc, err := blockchain.Open(memdb.NewBlockRepository(), nil)
	if err != nil {
		t.Fatalf("opening blockchain: %s", err)
	}

	actualHeader, err := bc.GetBestHeader()
	if err != nil {
		t.Fatalf("getting best header with blockchain: %s", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	server := rpc.NewServer(lis)
	server.RegisterDiscoveryServer(NewDiscoveryServer(bc))
	server.Serve()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTimeout(500*time.Millisecond), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewDiscoveryClient(conn)

	hash, err := client.Version()

	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	if !hash.IsEqual(actualHeader.Hash) {
		t.Errorf("headers don't match. got - %s, wanted - %s", hash, actualHeader.Hash)
	}
}

func TestDiscoveryServer_Version_TargetIsOffline(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:0", grpc.WithTimeout(500*time.Millisecond), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewDiscoveryClient(conn)

	_, err = client.Version()

	expectedErrorMessage := "rpc error: code = Unavailable desc = grpc: the connection is unavailable"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("unexpected error message. got - %s, wanted: %s", err, expectedErrorMessage)
	}
}
