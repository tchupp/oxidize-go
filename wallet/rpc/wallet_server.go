package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
)

type walletServer struct {
	backend account.Engine
}

func NewWalletServer(backend account.Engine) WalletServiceServer {
	return &walletServer{backend: backend}
}

func RegisterWalletServer(s *rpc.Server, srv WalletServiceServer) {
	s.Register(&_WalletService_serviceDesc, srv)
}

func (s *walletServer) Account(ctx context.Context, req *AccountRequest) (*AccountResponse, error) {
	addresses, err := mapAddresses(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	accounts, err := s.findAccounts(addresses)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	return &AccountResponse{Accounts: mapAccountsForResponse(accounts)}, nil
}

func mapAddresses(req *AccountRequest) ([]*identity.Address, error) {
	var result *multierror.Error
	var addresses []*identity.Address

	for _, addr := range req.Addresses {
		address, err := identity.DeserializeAddress(addr)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		addresses = append(addresses, address)
	}

	return addresses, result.ErrorOrNil()
}

func (s *walletServer) findAccounts(addresses []*identity.Address) ([]*account.Account, error) {
	var result *multierror.Error
	var accounts []*account.Account
	for _, address := range addresses {
		acc, err := s.backend.Account(address)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		accounts = append(accounts, acc)
	}
	return accounts, result.ErrorOrNil()
}

func mapAccountsForResponse(accounts []*account.Account) []*Account {
	var res []*Account
	for _, acc := range accounts {
		var txs []*Transaction
		for _, tx := range acc.Transactions() {
			txs = append(txs, &Transaction{
				Amount:   proto.Uint64(tx.Amount()),
				Spender:  proto.String(tx.Spender().Serialize()),
				Receiver: proto.String(tx.Receiver().Serialize()),
			})
		}

		res = append(res, &Account{
			Address:      proto.String(acc.Address().Serialize()),
			Spendable:    proto.Uint64(acc.Spendable()),
			Transactions: txs,
		})
	}
	return res
}
