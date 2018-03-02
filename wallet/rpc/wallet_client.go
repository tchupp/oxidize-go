package rpc

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/identity"
)

type WalletClient interface {
	Account([]*identity.Address) ([]*account.Account, error)
}

type walletClient struct {
	client WalletServiceClient
}

func NewWalletClient(conn *grpc.ClientConn) WalletClient {
	client := NewWalletServiceClient(conn)

	return &walletClient{client: client}
}

func (c *walletClient) Account(addresses []*identity.Address) ([]*account.Account, error) {
	var addrs []string
	for _, addr := range addresses {
		addrs = append(addrs, addr.Serialize())
	}

	response, err := c.client.Account(context.Background(), &AccountRequest{Addresses: addrs})
	if err != nil {
		return nil, err
	}

	accounts, err := mapAccountsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func mapAccountsFromResponse(response *AccountResponse) ([]*account.Account, error) {
	var result *multierror.Error
	var accounts []*account.Account
	for _, acc := range response.Accounts {
		address, err := identity.DeserializeAddress(acc.GetAddress())
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("deserializing account address '%s': %s", acc.GetAddress(), err))
			continue
		}

		txs, r := mapTransactionsFromAccount(acc)
		if r.ErrorOrNil() != nil {
			result = multierror.Append(result, r.WrappedErrors()...)
			continue
		}

		accounts = append(accounts, account.NewAccount(
			address,
			acc.GetSpendable(),
			txs,
		))
	}
	return accounts, result.ErrorOrNil()
}

func mapTransactionsFromAccount(acc *Account) ([]*account.Transaction, *multierror.Error) {
	var result *multierror.Error

	var txs []*account.Transaction
	for _, tx := range acc.Transactions {
		receiver, err := identity.DeserializeAddress(tx.GetReceiver())
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("deserializing receiver address '%s': %s", tx.GetReceiver(), err))
			continue
		}

		spender, err := identity.DeserializeAddress(tx.GetSpender())
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("deserializing spender address '%s': %s", tx.GetSpender(), err))
			continue
		}

		txs = append(txs, account.NewTransaction(tx.GetAmount(), spender, receiver))
	}
	return txs, result
}
