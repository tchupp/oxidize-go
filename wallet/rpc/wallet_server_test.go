package rpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/account/testdata"
	"github.com/tclchiam/oxidize-go/blockchain/engine/txsigning"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/server/rpc"
)

func TestWalletServer_Send(t *testing.T) {
	verifyUnspentOutputs := func(t *testing.T, client WalletClient, id *identity.Identity, index uint32, balance uint64) *UnspentOutputRef {
		unspentOutputs, err := client.UnspentOutputs([]*identity.Address{id.Address()})
		require.NoError(t, err, "error getting unspent outputs")
		require.Len(t, unspentOutputs, int(index+1), "unexpected number of outputs")

		unspentOutputRef := unspentOutputs[index]

		if !unspentOutputRef.Address.IsEqual(id.Address()) {
			t.Errorf("unexpected unpsentOutput address. Got - %s, wanted - %s", unspentOutputRef.Address, id.Address())
		}
		expectedOutput := &entity.Output{Index: index, Value: balance, PublicKeyHash: id.PublicKey().Hash()}
		if !unspentOutputRef.Output.IsEqual(expectedOutput) {
			t.Errorf("unexpected unpsentOutput output. Got - %s, wanted - %s", unspentOutputRef.Output, expectedOutput)
		}

		return unspentOutputRef
	}
	getAccounts := func(t *testing.T, client WalletClient, addresses []*identity.Address) map[string]*account.Account {
		accounts, err := client.Accounts(addresses)
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

	spender := identity.RandomIdentity()
	receiver := identity.RandomIdentity()
	addresses := []*identity.Address{spender.Address(), receiver.Address()}

	server, client := setup(t, spender)
	server.Serve()

	accounts := getAccounts(t, client, addresses)
	verifyBalance(t, addresses[0], 10, accounts[addresses[0].Serialize()])
	verifyBalance(t, addresses[1], 00, accounts[addresses[1].Serialize()])

	unspentOutputRef := verifyUnspentOutputs(t, client, spender, 0, 10)

	outputs := []*entity.Output{
		entity.NewOutput(3, spender.Address()),
		entity.NewOutput(7, receiver.Address()),
	}
	expenseTx := buildExpenseTx(unspentOutputRef, outputs, spender)
	assert.NoError(t, client.ProposeTransaction(expenseTx))

	accounts = getAccounts(t, client, addresses)
	verifyBalance(t, addresses[0], 3, accounts[addresses[0].Serialize()])
	verifyBalance(t, addresses[1], 7, accounts[addresses[1].Serialize()])

	verifyUnspentOutputs(t, client, receiver, 0, 7)

	assert.NoError(t, server.Close())
}

func buildExpenseTx(unspentOutputRef *UnspentOutputRef, outputs []*entity.Output, spender *identity.Identity) *entity.Transaction {
	unsignedInput := entity.NewUnsignedInput(unspentOutputRef.TxId, unspentOutputRef.Output, spender.PublicKey())
	signature := txsigning.GenerateSignature(unsignedInput, outputs, spender, encoding.TransactionProtoEncoder())

	signedInputs := []*entity.SignedInput{entity.NewSignedInput(unsignedInput, signature)}
	return entity.NewTx(signedInputs, outputs, encoding.TransactionProtoEncoder())
}

func setup(t *testing.T, owner *identity.Identity) (*rpc.Server, WalletClient) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "starting listener")

	engine := testdata.NewAccountEngineBuilder(t).
		Build().
		AddBalance(owner.Address(), 10)

	server := rpc.NewServer(lis)
	RegisterWalletServer(server, NewWalletServer(engine))

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err, "dialing server")

	return server, NewWalletClient(conn)
}
