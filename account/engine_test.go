package account

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/testdata"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
)

func Test_engine_Balance(t *testing.T) {
	type testRun struct {
		name    string
		bc      blockchain.Blockchain
		address *identity.Address
		want    *Account
		wantErr bool
	}

	brokeAddress := identity.RandomIdentity().Address()
	brokeAccount := testRun{
		name:    "broke account",
		bc:      testdata.NewBlockchainBuilder(t).Build(),
		want:    &Account{Address: brokeAddress, Spendable: 0},
		address: brokeAddress,
	}

	richAddress := identity.RandomIdentity().Address()
	richAccount := testRun{
		name: "rich account",
		bc: testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(richAddress, 1000),
		want:    &Account{Address: richAddress, Spendable: 1000},
		address: richAddress,
	}

	for _, tt := range []testRun{brokeAccount, richAccount} {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine(tt.bc)

			got, err := engine.Balance(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("engine.Balance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.IsEqual(tt.want) {
				t.Errorf("engine.Balance() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, engine.Close())
		})
	}
}

func Test_engine_Transactions(t *testing.T) {
	t.Run("engine.Transactions() - none", func(t *testing.T) {
		engine := NewEngine(testdata.NewBlockchainBuilder(t).Build())

		brokeIdentity := identity.RandomIdentity()

		got, err := engine.Transactions(brokeIdentity.Address())
		if err != nil {
			t.Errorf("engine.Transactions() error = %v, wantErr %v", err, false)
			return
		}
		if !reflect.DeepEqual(got, Transactions(nil)) {
			t.Errorf("engine.Transactions() = %v, want %v", got, Transactions(nil))
		}

		assert.NoError(t, engine.Close())
	})

	t.Run("engine.Transactions() - spending", func(t *testing.T) {
		spendingIdentity := identity.RandomIdentity()
		receivingIdentity := identity.RandomIdentity()

		engine := NewEngine(testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(spendingIdentity.Address(), 10))

		err := engine.Send(spendingIdentity, receivingIdentity.Address(), 10)
		if err != nil {
			t.Errorf("engine.Send() error = %v, wantErr %v", err, false)
			return
		}

		got, err := engine.Transactions(spendingIdentity.Address())
		if err != nil {
			t.Errorf("engine.Transactions() error = %v, wantErr %v", err, false)
			return
		}
		expectedTxs := Transactions{
			{amount: 10, spender: nil, receiver: spendingIdentity.Address()},
			{amount: 10, spender: spendingIdentity.Address(), receiver: receivingIdentity.Address()},
		}
		if !reflect.DeepEqual(got, expectedTxs) {
			t.Errorf("engine.Transactions() = %v, want %v", got, expectedTxs)
		}

		assert.NoError(t, engine.Close())
	})

	t.Run("engine.Transactions() - receiving", func(t *testing.T) {
		spendingIdentity := identity.RandomIdentity()
		receivingIdentity := identity.RandomIdentity()

		engine := NewEngine(testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(spendingIdentity.Address(), 10))

		err := engine.Send(spendingIdentity, receivingIdentity.Address(), 10)
		if err != nil {
			t.Errorf("engine.Send() error = %v, wantErr %v", err, false)
			return
		}

		got, err := engine.Transactions(receivingIdentity.Address())
		if err != nil {
			t.Errorf("engine.Transactions() error = %v, wantErr %v", err, false)
			return
		}
		expectedTxs := Transactions{{amount: 10, spender: spendingIdentity.Address(), receiver: receivingIdentity.Address()}}
		if !reflect.DeepEqual(got, expectedTxs) {
			t.Errorf("engine.Transactions() = %v, want %v", got, expectedTxs)
		}

		assert.NoError(t, engine.Close())
	})

	t.Run("engine.Transactions() - reward", func(t *testing.T) {
		spendingIdentity := identity.RandomIdentity()
		receivingIdentity := identity.RandomIdentity()

		bc := testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(spendingIdentity.Address(), 10)

		rewardTx := entity.NewRewardTx(receivingIdentity.Address(), encoding.TransactionProtoEncoder())
		block, err := bc.MineBlock(entity.Transactions{rewardTx})
		if bc.SaveBlock(block); err != nil {
			t.Errorf("engine.Send() error = %v, wantErr %v", err, false)
			return
		}

		engine := NewEngine(bc)
		got, err := engine.Transactions(receivingIdentity.Address())
		if err != nil {
			t.Errorf("engine.Transactions() error = %v, wantErr %v", err, false)
			return
		}
		expectedTxs := Transactions{{amount: 10, spender: nil, receiver: receivingIdentity.Address()}}
		if !reflect.DeepEqual(got, expectedTxs) {
			t.Errorf("engine.Transactions() = %v, want %v", got, expectedTxs)
		}

		assert.NoError(t, engine.Close())
	})
}

func Test_engine_Send(t *testing.T) {
	spender := identity.RandomIdentity()
	receiver := identity.RandomIdentity().Address()

	type args struct {
		spender  *identity.Identity
		receiver *identity.Address
		expense  uint64
	}
	type testState struct {
		spenderBalance  uint64
		receiverBalance uint64
	}
	tests := []struct {
		name    string
		bc      blockchain.Blockchain
		args    args
		wantErr bool
		before  testState
		after   testState
	}{
		{
			name:   "sending 0",
			bc:     testdata.NewBlockchainBuilder(t).Build(),
			args:   args{spender: spender, receiver: receiver, expense: 0},
			before: testState{spenderBalance: 0, receiverBalance: 0},
			after:  testState{spenderBalance: 0, receiverBalance: 0},
		},
		{
			name:    "over spending",
			bc:      testdata.NewBlockchainBuilder(t).Build(),
			args:    args{spender: spender, receiver: receiver, expense: 10},
			wantErr: true,
			before:  testState{spenderBalance: 0, receiverBalance: 0},
			after:   testState{spenderBalance: 0, receiverBalance: 0},
		},
		{
			name: "under spending",
			bc: testdata.NewBlockchainBuilder(t).
				Build().
				AddBalance(spender.Address(), 20),
			args:   args{spender: spender, receiver: receiver, expense: 10},
			before: testState{spenderBalance: 20, receiverBalance: 0},
			after:  testState{spenderBalance: 10, receiverBalance: 10},
		},
		{
			name: "exact spending",
			bc: testdata.NewBlockchainBuilder(t).
				Build().
				AddBalance(spender.Address(), 10),
			args:   args{spender: spender, receiver: receiver, expense: 10},
			before: testState{spenderBalance: 10, receiverBalance: 0},
			after:  testState{spenderBalance: 0, receiverBalance: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine(tt.bc)

			if account, err := engine.Balance(tt.args.spender.Address()); err != nil {
				if account.Spendable != tt.before.spenderBalance {
					t.Errorf("spender does not have expected before balance. want - %d, got - %d", tt.before.spenderBalance, account.Spendable)
				}
			}
			if account, err := engine.Balance(tt.args.receiver); err != nil {
				if account.Spendable != tt.before.receiverBalance {
					t.Errorf("receiver does not have expected before balance. want - %d, got - %d", tt.before.receiverBalance, account.Spendable)
				}
			}

			if err := engine.Send(tt.args.spender, tt.args.receiver, tt.args.expense); (err != nil) != tt.wantErr {
				t.Errorf("engine.Send() error = %v, wantErr %v", err, tt.wantErr)
			}

			if account, err := engine.Balance(tt.args.spender.Address()); err != nil {
				if account.Spendable != tt.after.spenderBalance {
					t.Errorf("spender does not have expected after balance. want - %d, got - %d", tt.after.spenderBalance, account.Spendable)
				}
			}
			if account, err := engine.Balance(tt.args.receiver); err != nil {
				if account.Spendable != tt.after.receiverBalance {
					t.Errorf("receiver does not have expected after balance. want - %d, got - %d", tt.after.receiverBalance, account.Spendable)
				}
			}

			assert.NoError(t, engine.Close())
		})
	}
}
