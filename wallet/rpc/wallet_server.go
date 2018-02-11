package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/rpc"
)

type walletBackend interface {
	Balance(*identity.Address) (*account.Account, error)
}

type walletServer struct {
	backend walletBackend
}

func NewWalletServer(backend walletBackend) WalletServiceServer {
	return &walletServer{backend: backend}
}

func RegisterWalletServer(s *rpc.Server, srv WalletServiceServer) {
	s.Register(&_WalletService_serviceDesc, srv)
}

func (s *walletServer) Balance(ctx context.Context, req *BalanceRequest) (*BalanceResponse, error) {
	addresses, err := mapAddresses(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	accounts, err := s.findBalances(addresses)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	return &BalanceResponse{Accounts: mapAccountsForResponse(accounts)}, nil
}

func mapAddresses(req *BalanceRequest) ([]*identity.Address, error) {
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

func (s *walletServer) findBalances(addresses []*identity.Address) ([]*account.Account, error) {
	var result *multierror.Error
	var accounts []*account.Account
	for _, address := range addresses {
		acc, err := s.backend.Balance(address)
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
		res = append(res, &Account{
			Address:   proto.String(acc.Address.Serialize()),
			Total:     proto.Uint64(0),
			Spendable: proto.Uint64(acc.Spendable),
			Reward:    proto.Uint64(0),
		})
	}
	return res
}
