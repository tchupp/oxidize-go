package rpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage/memdb"
)

func TestWalletServer_Balance(t *testing.T) {
	owner := identity.RandomIdentity()
	receiver := identity.RandomIdentity()

	engine := setupAccountEngine(t, owner)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	server := rpc.NewServer(lis)
	RegisterWalletServer(server, NewWalletServer(engine))
	server.Serve()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	client := NewWalletClient(conn)
	accounts, err := client.Balance([]*identity.Address{owner.Address()})
	if err != nil {
		t.Fatalf("error getting balance: %s", err)
	}

	if len(accounts) != 1 {
		t.Errorf("unexpected accounts. got - %d, wanted - %d.", len(accounts), 1)
	}

	actualOwnerAccount := accounts[0]
	expectedOwnerAccount := &account.Account{Address: owner.Address(), Spendable: 10}

	if !actualOwnerAccount.IsEqual(expectedOwnerAccount) {
		t.Errorf("initial owner account incorrect. got - %s, wanted - %s", actualOwnerAccount, expectedOwnerAccount)
	}

	err = engine.Send(owner, receiver.Address(), 7)
	if err != nil {
		t.Fatalf("failed to send: %s", err)
	}

	accounts, err = client.Balance([]*identity.Address{owner.Address(), receiver.Address()})
	if err != nil {
		t.Fatalf("error getting balance: %s", err)
	}

	if len(accounts) != 2 {
		t.Errorf("unexpected accounts. got - %d, wanted - %d.", len(accounts), 1)
	}

	actualOwnerAccount = accounts[0]
	expectedOwnerAccount = &account.Account{Address: owner.Address(), Spendable: 13}

	if !actualOwnerAccount.IsEqual(expectedOwnerAccount) {
		t.Errorf("initial owner account incorrect. got - %s, wanted - %s", actualOwnerAccount, expectedOwnerAccount)
	}

	actualReceiverAccount := accounts[1]
	expectedReceiverAccount := &account.Account{Address: receiver.Address(), Spendable: 7}

	if !actualReceiverAccount.IsEqual(expectedReceiverAccount) {
		t.Errorf("initial owner account incorrect. got - %s, wanted - %s", actualReceiverAccount, expectedReceiverAccount)
	}
}

func setupAccountEngine(t *testing.T, owner *identity.Identity) account.Engine {
	miner := proofofwork.NewDefaultMiner(owner.Address())

	repository := memdb.NewChainRepository()
	genesisBlock := miner.MineBlock(&entity.GenesisParentHeader, entity.Transactions{})
	if err := repository.SaveBlock(genesisBlock); err != nil {
		t.Fatalf("failed to save genesis")
	}

	bc, err := blockchain.Open(repository, miner)
	if err != nil {
		t.Fatalf("failed to open blockchain")
	}

	return account.NewEngine(bc)
}
