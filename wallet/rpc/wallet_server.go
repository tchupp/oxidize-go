package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/server/rpc"
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
	addresses, err := mapAddresses(req.Addresses)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	accounts, err := findAccounts(s.backend, addresses)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	return &AccountResponse{Accounts: mapAccountsToResponse(accounts)}, nil
}

func mapAddresses(addrs []string) ([]*identity.Address, error) {
	var result *multierror.Error
	var addresses []*identity.Address

	for _, addr := range addrs {
		address, err := identity.DeserializeAddress(addr)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		addresses = append(addresses, address)
	}

	return addresses, result.ErrorOrNil()
}

func findAccounts(engine account.Engine, addresses []*identity.Address) ([]*account.Account, error) {
	var result *multierror.Error
	var accounts []*account.Account
	for _, address := range addresses {
		acc, err := engine.Account(address)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		accounts = append(accounts, acc)
	}
	return accounts, result.ErrorOrNil()
}

func mapAccountsToResponse(accounts []*account.Account) []*Account {
	var res []*Account
	for _, acc := range accounts {
		var txs []*encoding.Transaction
		for _, transaction := range acc.Transactions() {
			txs = append(txs, encoding.ToWireTransaction(transaction))
		}

		res = append(res, &Account{
			Address:      proto.String(acc.Address().Serialize()),
			Spendable:    proto.Uint64(acc.Spendable()),
			Transactions: txs,
		})
	}
	return res
}

func (s *walletServer) UnspentOutputs(ctx context.Context, req *UnspentOutputsRequest) (*UnspentOutputsResponse, error) {
	addresses, err := mapAddresses(req.Addresses)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested addresses were invalid: %s", err)
	}

	return mapOutputsToResponse(s.backend, addresses)
}

func mapOutputsToResponse(engine account.Engine, addresses []*identity.Address) (*UnspentOutputsResponse, error) {
	res := &UnspentOutputsResponse{}

	var result *multierror.Error
	for _, address := range addresses {
		spendableOutputs, err := engine.SpendableOutputs(address)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		spendableOutputs.ForEach(func(txId *entity.Hash, output *entity.Output) {
			res.Outputs = append(res.Outputs, &WireUnspentOutput{
				Address: proto.String(address.Serialize()),
				TxHash:  txId.Slice(),
				Output:  encoding.ToWireOutput(output),
			})
		})
	}
	return res, result.ErrorOrNil()
}

func (s *walletServer) ProposeTransaction(ctx context.Context, req *ProposeTransactionRequest) (*ProposeTransactionResponse, error) {
	transaction, err := encoding.FromWireTransaction(req.GetTransaction())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested transaction was invalid: %s", err)
	}

	if err := s.backend.ProposeTransaction(transaction); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to propose transaction: %s", err)
	}

	return &ProposeTransactionResponse{}, nil
}
