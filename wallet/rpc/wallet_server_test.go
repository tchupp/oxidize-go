package rpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/account/testdata"
	"github.com/tclchiam/oxidize-go/blockchain/engine/txsigning"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
)

func TestWalletServer_Balance(t *testing.T) {
	verifyBalance := func(t *testing.T, id *identity.Identity, balance uint64, actualAccount *account.Account) {
		expectedAccount := account.NewAccount(id.Address(), balance, nil)

		if !actualAccount.IsEqual(expectedAccount) {
			t.Errorf("unexpected account. got - %s, wanted - %s", actualAccount, expectedAccount)
		}
	}

	owner := identity.RandomIdentity()
	receiver := identity.RandomIdentity()

	server, engine, client := setup(t, owner)
	server.Serve()

	accounts, err := client.Account([]*identity.Address{owner.Address(), receiver.Address()})
	if err != nil {
		t.Fatalf("error getting balance: %s", err)
	}
	if !assert.Len(t, accounts, 2) {
		t.FailNow()
	}

	verifyBalance(t, owner, 10, accounts[0])
	verifyBalance(t, receiver, 0, accounts[1])

	err = engine.Send(owner, receiver.Address(), 7)
	if err != nil {
		t.Fatalf("failed to send: %s", err)
	}

	accounts, err = client.Account([]*identity.Address{owner.Address(), receiver.Address()})
	if err != nil {
		t.Fatalf("error getting balance: %s", err)
	}
	if !assert.Len(t, accounts, 2) {
		t.FailNow()
	}

	verifyBalance(t, owner, 3, accounts[0])
	verifyBalance(t, receiver, 7, accounts[1])

	assert.NoError(t, server.Close())
}

func TestWalletServer_Send(t *testing.T) {
	verifyUnspentOutputs := func(t *testing.T, client WalletClient, id *identity.Identity) *UnspentOutputRef {
		unspentOutputs, err := client.UnspentOutputs([]*identity.Address{id.Address()})
		if err != nil {
			t.Fatalf("error getting unspent outputs: %s", err)
		}
		if len(unspentOutputs) != 1 {
			t.Fatalf("unexpected number of outputs. Got - %d, wanted - %d", len(unspentOutputs), 1)
		}

		unspentOutputRef := unspentOutputs[0]

		if !unspentOutputRef.Address.IsEqual(id.Address()) {
			t.Errorf("unexpected unpsentOutput address. Got - %s, wanted - %s", unspentOutputRef.Address, id.Address())
		}
		expectedOutput := &entity.Output{Index: 0, Value: 10, PublicKeyHash: id.PublicKey().Hash()}
		if !unspentOutputRef.Output.IsEqual(expectedOutput) {
			t.Errorf("unexpected unpsentOutput output. Got - %s, wanted - %s", unspentOutputRef.Output, expectedOutput)
		}

		return unspentOutputRef
	}

	spender := identity.RandomIdentity()
	receiver := identity.RandomIdentity()

	server, _, client := setup(t, spender)
	server.Serve()

	unspentOutputRef := verifyUnspentOutputs(t, client, spender)

	expenseTx := buildExpenseTx(unspentOutputRef, []*entity.Output{entity.NewOutput(10, receiver.Address())}, spender)
	assert.NoError(t, client.ProposeTransaction(expenseTx))

	verifyUnspentOutputs(t, client, receiver)

	assert.NoError(t, server.Close())
}

func buildExpenseTx(unspentOutputRef *UnspentOutputRef, outputs []*entity.Output, spender *identity.Identity) *entity.Transaction {
	unsignedInput := entity.NewUnsignedInput(unspentOutputRef.TxId, unspentOutputRef.Output, spender.PublicKey())
	signature := txsigning.GenerateSignature(unsignedInput, outputs, spender, encoding.TransactionProtoEncoder())
	signedInputs := []*entity.SignedInput{entity.NewSignedInput(unsignedInput, signature)}
	return entity.NewTx(signedInputs, outputs, encoding.TransactionProtoEncoder())
}

func setup(t *testing.T, owner *identity.Identity) (*rpc.Server, *testdata.TestAccountEngine, WalletClient) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting listener: %s", err)
	}

	engine := testdata.NewAccountEngineBuilder(t).
		Build().
		AddBalance(owner.Address(), 10)

	server := rpc.NewServer(lis)
	RegisterWalletServer(server, NewWalletServer(engine))

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dialing server: %s", err)
	}

	return server, engine, NewWalletClient(conn)
}
