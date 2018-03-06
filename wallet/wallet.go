package wallet

import (
	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/wallet/rpc"
)

type UnspentOutput struct {
	Identity *identity.Identity
	TxId     *entity.Hash
	Output   *entity.Output
}

type Wallet interface {
	Identities() (identity.Identities, error)
	NewIdentity() (*identity.Identity, error)
	Accounts() ([]*account.Account, error)
	Send(*identity.Address, *identity.Address, uint64) error
}

type wallet struct {
	*KeyStore
	rpc.WalletClient
}

func NewWallet(store *KeyStore, client rpc.WalletClient) Wallet {
	return &wallet{KeyStore: store, WalletClient: client}
}

func (w *wallet) NewIdentity() (*identity.Identity, error) {
	newIdentity := identity.RandomIdentity()
	err := w.KeyStore.SaveIdentity(newIdentity)
	if err != nil {
		return nil, err
	}

	return newIdentity, nil
}

func (w *wallet) Accounts() ([]*account.Account, error) {
	addresses, err := addresses(w.Identities())
	if err != nil {
		return nil, err
	}

	return w.WalletClient.Accounts(addresses)
}

func (w *wallet) Send(receiver, payback *identity.Address, amount uint64) error {
	addresses, err := addresses(w.Identities())
	if err != nil {
		return err
	}

	outputs, err := w.WalletClient.UnspentOutputs(addresses)
	if err != nil {
		return err
	}

	unspentOutputs, err := mapToUnspentOutputs(w.KeyStore, outputs)

	transaction, err := buildExpenseTransaction(unspentOutputs, receiver, payback, amount)
	if err != nil {
		return err
	}

	return w.WalletClient.ProposeTransaction(transaction)
}

func mapToUnspentOutputs(keyStore *KeyStore, outputs []*rpc.UnspentOutputRef) (unspentOutputs []*UnspentOutput, err error) {
	var result *multierror.Error
	for _, output := range outputs {
		id, err := keyStore.Identity(output.Address.Serialize())
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		unspentOutputs = append(unspentOutputs, &UnspentOutput{
			Identity: id,
			TxId:     output.TxId,
			Output:   output.Output,
		})
	}

	return unspentOutputs, result.ErrorOrNil()
}

func addresses(ids []*identity.Identity, err error) ([]*identity.Address, error) {
	if err != nil {
		return nil, err
	}

	var addrs []*identity.Address
	for _, id := range ids {
		addrs = append(addrs, id.Address())
	}

	return addrs, nil
}
