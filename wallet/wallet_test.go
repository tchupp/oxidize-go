package wallet

import (
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/account/testdata"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
	walletrpc "github.com/tclchiam/oxidize-go/wallet/rpc"
)

func TestWalletServer_Send(t *testing.T) {
	getAccounts := func(t *testing.T, wallet Wallet, addresses []*identity.Address) map[string]*account.Account {
		accounts, err := wallet.Accounts()
		require.NoError(t, err, "error getting accounts")
		require.Len(t, accounts, len(addresses))

		accts := make(map[string]*account.Account, 0)
		for _, acct := range accounts {
			accts[acct.Address().Serialize()] = acct
		}
		return accts
	}
	verifyBalance := func(t *testing.T, address *identity.Address, balance uint64, actualAccount *account.Account) {
		expectedAccount := account.NewAccount(address, balance, nil)

		if !actualAccount.IsEqual(expectedAccount) {
			t.Errorf("unexpected account. got - %s, wanted - %s", actualAccount, expectedAccount)
		}
	}

	keyStore := NewKeyStore(makeKeystoreDir())
	defer os.RemoveAll(keyStore.path)

	spender1 := identity.RandomIdentity()
	spender2 := identity.RandomIdentity()
	receiver := identity.RandomIdentity()

	require.NoError(t, keyStore.SaveIdentity(spender1), "error saving spender1 id")
	require.NoError(t, keyStore.SaveIdentity(spender2), "error saving spender2 id")
	require.NoError(t, keyStore.SaveIdentity(receiver), "error saving receiver id")

	addresses := []*identity.Address{spender1.Address(), spender2.Address(), receiver.Address()}

	type args struct {
		expense int64
	}

	tests := []struct {
		name    string
		engine  account.Engine
		args    args
		wantErr bool
		before  map[string]uint64
		after   map[string]uint64
	}{
		{
			name:    "sending 0",
			engine:  testdata.NewAccountEngineBuilder(t).Build(),
			args:    args{expense: 0},
			wantErr: true,
			before: map[string]uint64{
				spender1.Address().Serialize(): 0,
				spender2.Address().Serialize(): 0,
				receiver.Address().Serialize(): 0,
			},
			after: map[string]uint64{
				spender1.Address().Serialize(): 0,
				spender2.Address().Serialize(): 0,
				receiver.Address().Serialize(): 0,
			},
		},
		{
			name:    "over spending",
			engine:  testdata.NewAccountEngineBuilder(t).Build(),
			args:    args{expense: 10},
			wantErr: true,
			before: map[string]uint64{
				spender1.Address().Serialize(): 0,
				spender2.Address().Serialize(): 0,
				receiver.Address().Serialize(): 0,
			},
			after: map[string]uint64{
				spender1.Address().Serialize(): 0,
				spender2.Address().Serialize(): 0,
				receiver.Address().Serialize(): 0,
			},
		},
		{
			name: "under spending",
			engine: testdata.NewAccountEngineBuilder(t).
				Build().
				AddBalance(spender1.Address(), 20).
				AddBalance(spender2.Address(), 20),
			args: args{expense: 30},
			before: map[string]uint64{
				spender1.Address().Serialize(): 20,
				spender2.Address().Serialize(): 20,
				receiver.Address().Serialize(): 0,
			},
			after: map[string]uint64{
				spender1.Address().Serialize(): 10,
				spender2.Address().Serialize(): 0,
				receiver.Address().Serialize(): 30,
			},
		},
		{
			name: "exact spending",
			engine: testdata.NewAccountEngineBuilder(t).
				Build().
				AddBalance(spender1.Address(), 20).
				AddBalance(spender2.Address(), 20),
			args: args{expense: 40},
			before: map[string]uint64{
				spender1.Address().Serialize(): 20,
				spender2.Address().Serialize(): 20,
				receiver.Address().Serialize(): 0,
			},
			after: map[string]uint64{
				spender1.Address().Serialize(): 0,
				spender2.Address().Serialize(): 0,
				receiver.Address().Serialize(): 40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupServer(t, tt.engine)
			wallet := setupWallet(t, keyStore, server.Addr())
			defer func() {
				assert.NoError(t, server.Close())
				assert.NoError(t, tt.engine.Close())
			}()

			accounts := getAccounts(t, wallet, addresses)
			verifyBalance(t, addresses[0], tt.before[addresses[0].Serialize()], accounts[addresses[0].Serialize()])
			verifyBalance(t, addresses[1], tt.before[addresses[1].Serialize()], accounts[addresses[1].Serialize()])
			verifyBalance(t, addresses[2], tt.before[addresses[2].Serialize()], accounts[addresses[2].Serialize()])

			if err := wallet.Send(receiver.Address(), spender1.Address(), tt.args.expense); (err != nil) != tt.wantErr {
				t.Errorf("wallet.Send(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			accounts = getAccounts(t, wallet, addresses)
			verifyBalance(t, addresses[0], tt.after[addresses[0].Serialize()], accounts[addresses[0].Serialize()])
			verifyBalance(t, addresses[1], tt.after[addresses[1].Serialize()], accounts[addresses[1].Serialize()])
			verifyBalance(t, addresses[2], tt.after[addresses[2].Serialize()], accounts[addresses[2].Serialize()])
		})
	}
}

func setupServer(t *testing.T, engine account.Engine) *rpc.Server {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "starting listener")

	server := rpc.NewServer(lis)
	walletrpc.RegisterWalletServer(server, walletrpc.NewWalletServer(engine))
	server.Serve()

	return server
}

func setupWallet(t *testing.T, keyStore *KeyStore, addr net.Addr) Wallet {
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	require.NoError(t, err, "dialing server")

	return NewWallet(keyStore, walletrpc.NewWalletClient(conn))
}
