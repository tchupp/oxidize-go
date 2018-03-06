package account

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/testdata"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
)

func Test_engine_Account(t *testing.T) {
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
		want:    &Account{address: brokeAddress, spendable: 0},
		address: brokeAddress,
	}

	richAddress := identity.RandomIdentity().Address()
	richAccount := testRun{
		name: "rich account",
		bc: testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(richAddress, 1000),
		want:    &Account{address: richAddress, spendable: 1000},
		address: richAddress,
	}

	for _, tt := range []testRun{brokeAccount, richAccount} {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine(tt.bc)

			got, err := engine.Account(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("engine.Account() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.IsEqual(tt.want) {
				t.Errorf("engine.Account() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, engine.Close())
		})
	}
}

func Test_engine_Transactions(t *testing.T) {
	t.Run("engine.Account() - none", func(t *testing.T) {
		engine := NewEngine(testdata.NewBlockchainBuilder(t).Build())

		brokeIdentity := identity.RandomIdentity()

		got, err := engine.Account(brokeIdentity.Address())
		if err != nil {
			t.Errorf("engine.Account() error = %v, wantErr %v", err, false)
			return
		}
		if !reflect.DeepEqual(got.Transactions(), Transactions(nil)) {
			t.Errorf("engine.Account() = %v, want %v", got.Transactions(), Transactions(nil))
		}

		assert.NoError(t, engine.Close())
	})

	t.Run("engine.Account() - spending", func(t *testing.T) {
		spendingIdentity := identity.RandomIdentity()

		engine := NewEngine(testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(spendingIdentity.Address(), 10))

		got, err := engine.Account(spendingIdentity.Address())
		require.NoError(t, err, "engine.Account()")

		expectedTxs := Transactions{
			{amount: 10, spender: nil, receiver: spendingIdentity.Address()},
		}
		if !reflect.DeepEqual(got.Transactions(), expectedTxs) {
			t.Errorf("engine.Account() = %v, want %v", got.Transactions(), expectedTxs)
		}

		assert.NoError(t, engine.Close())
	})

	t.Run("engine.Account() - reward", func(t *testing.T) {
		spendingIdentity := identity.RandomIdentity()
		receivingIdentity := identity.RandomIdentity()

		bc := testdata.NewBlockchainBuilder(t).
			Build().
			AddBalance(spendingIdentity.Address(), 10)

		rewardTx := entity.NewRewardTx(receivingIdentity.Address(), encoding.TransactionProtoEncoder())
		block, err := bc.MineBlock(entity.Transactions{rewardTx})
		require.NoError(t, err, "mining block")

		err = bc.SaveBlock(block)
		require.NoError(t, err, "saving block")

		engine := NewEngine(bc)
		got, err := engine.Account(receivingIdentity.Address())
		require.NoError(t, err, "engine.Account()")

		expectedTxs := Transactions{{amount: 10, spender: nil, receiver: receivingIdentity.Address()}}
		if !reflect.DeepEqual(got.Transactions(), expectedTxs) {
			t.Errorf("engine.Account() = %v, want %v", got.Transactions(), expectedTxs)
		}

		assert.NoError(t, engine.Close())
	})
}
