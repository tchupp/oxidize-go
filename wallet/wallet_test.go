package wallet

import (
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/account/testdata"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
	walletrpc "github.com/tclchiam/oxidize-go/wallet/rpc"
)

func TestWalletServer_Send(t *testing.T) {
	getAccounts := func(t *testing.T, wallet Wallet, addresses []*identity.Address) []*account.Account {
		accounts, err := wallet.Accounts()
		if err != nil {
			t.Fatalf("error getting balance: %s", err)
		}
		if !assert.Len(t, accounts, len(addresses)) {
			t.FailNow()
		}

		return accounts
	}
	verifyBalance := func(t *testing.T, address *identity.Address, balance uint64, actualAccount *account.Account) {
		expectedAccount := account.NewAccount(address, balance, nil)

		if !actualAccount.IsEqual(expectedAccount) {
			t.Errorf("unexpected account. got - %s, wanted - %s", actualAccount, expectedAccount)
		}
	}

	wallet, engine, cleanup := setup(t)
	defer cleanup()

	spender, err := wallet.NewIdentity()
	assert.NoError(t, err, "error saving spender id")

	receiver, err := wallet.NewIdentity()
	assert.NoError(t, err, "error saving receiver id")

	engine.AddBalance(spender.Address(), 10)

	addresses := []*identity.Address{spender.Address(), receiver.Address()}

	accounts := getAccounts(t, wallet, addresses)
	verifyBalance(t, addresses[0], 10, accounts[0])
	verifyBalance(t, addresses[1], 0, accounts[1])

	err = wallet.Send(receiver.Address(), spender.Address(), 7)
	assert.NoError(t, err)

	accounts = getAccounts(t, wallet, addresses)
	verifyBalance(t, addresses[0], 3, accounts[0])
	verifyBalance(t, addresses[1], 7, accounts[1])
}

func setup(t *testing.T) (Wallet, *testdata.TestAccountEngine, func()) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	engine := testdata.NewAccountEngineBuilder(t).Build()

	server := rpc.NewServer(lis)
	walletrpc.RegisterWalletServer(server, walletrpc.NewWalletServer(engine))
	server.Serve()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	keyStore := NewKeyStore(makeKeystoreDir())

	cleanup := func() {
		assert.NoError(t, server.Close())
		os.RemoveAll(keyStore.path)
	}
	return NewWallet(keyStore, walletrpc.NewWalletClient(conn)), engine, cleanup
}
