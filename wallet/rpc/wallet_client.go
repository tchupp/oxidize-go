package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/identity"
)

type WalletClient interface {
	Balance([]*identity.Address) ([]*account.Account, error)
}

type walletClient struct {
	client WalletServiceClient
}

func NewWalletClient(conn *grpc.ClientConn) WalletClient {
	client := NewWalletServiceClient(conn)

	return &walletClient{client: client}
}

func (c *walletClient) Balance(addresses []*identity.Address) ([]*account.Account, error) {
	var addrs []string
	for _, addr := range addresses {
		addrs = append(addrs, addr.Serialize())
	}

	response, err := c.client.Balance(context.Background(), &BalanceRequest{Addresses: addrs})
	if err != nil {
		return nil, err
	}

	accounts, err := mapAccountsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func mapAccountsFromResponse(response *BalanceResponse) ([]*account.Account, error) {
	var result *multierror.Error
	var accounts []*account.Account
	for _, acc := range response.Accounts {
		address, err := identity.DeserializeAddress(acc.GetAddress())
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		accounts = append(accounts, &account.Account{
			Address:   address,
			Spendable: acc.GetSpendable(),
		})
	}
	return accounts, result.ErrorOrNil()
}
