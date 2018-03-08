// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wallet.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	wallet.proto

It has these top-level messages:
	Account
	WireUnspentOutput
	AccountRequest
	AccountResponse
	UnspentOutputsRequest
	UnspentOutputsResponse
	ProposeTransactionRequest
	ProposeTransactionResponse
*/
package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import entities "github.com/tclchiam/oxidize-go/encoding"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Account struct {
	Address          *string                 `protobuf:"bytes,1,req,name=address" json:"address,omitempty"`
	Spendable        *uint64                 `protobuf:"varint,3,req,name=spendable" json:"spendable,omitempty"`
	Transactions     []*entities.Transaction `protobuf:"bytes,5,rep,name=transactions" json:"transactions,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *Account) Reset()                    { *m = Account{} }
func (m *Account) String() string            { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()               {}
func (*Account) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Account) GetAddress() string {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return ""
}

func (m *Account) GetSpendable() uint64 {
	if m != nil && m.Spendable != nil {
		return *m.Spendable
	}
	return 0
}

func (m *Account) GetTransactions() []*entities.Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type WireUnspentOutput struct {
	Address          *string          `protobuf:"bytes,1,req,name=address" json:"address,omitempty"`
	TxHash           []byte           `protobuf:"bytes,2,req,name=txHash" json:"txHash,omitempty"`
	Output           *entities.Output `protobuf:"bytes,3,req,name=output" json:"output,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *WireUnspentOutput) Reset()                    { *m = WireUnspentOutput{} }
func (m *WireUnspentOutput) String() string            { return proto.CompactTextString(m) }
func (*WireUnspentOutput) ProtoMessage()               {}
func (*WireUnspentOutput) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *WireUnspentOutput) GetAddress() string {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return ""
}

func (m *WireUnspentOutput) GetTxHash() []byte {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *WireUnspentOutput) GetOutput() *entities.Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type AccountRequest struct {
	Addresses        []string `protobuf:"bytes,1,rep,name=addresses" json:"addresses,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *AccountRequest) Reset()                    { *m = AccountRequest{} }
func (m *AccountRequest) String() string            { return proto.CompactTextString(m) }
func (*AccountRequest) ProtoMessage()               {}
func (*AccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AccountRequest) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type AccountResponse struct {
	Accounts         []*Account `protobuf:"bytes,1,rep,name=accounts" json:"accounts,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *AccountResponse) Reset()                    { *m = AccountResponse{} }
func (m *AccountResponse) String() string            { return proto.CompactTextString(m) }
func (*AccountResponse) ProtoMessage()               {}
func (*AccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AccountResponse) GetAccounts() []*Account {
	if m != nil {
		return m.Accounts
	}
	return nil
}

type UnspentOutputsRequest struct {
	Addresses        []string `protobuf:"bytes,1,rep,name=addresses" json:"addresses,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *UnspentOutputsRequest) Reset()                    { *m = UnspentOutputsRequest{} }
func (m *UnspentOutputsRequest) String() string            { return proto.CompactTextString(m) }
func (*UnspentOutputsRequest) ProtoMessage()               {}
func (*UnspentOutputsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *UnspentOutputsRequest) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type UnspentOutputsResponse struct {
	Outputs          []*WireUnspentOutput `protobuf:"bytes,1,rep,name=outputs" json:"outputs,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *UnspentOutputsResponse) Reset()                    { *m = UnspentOutputsResponse{} }
func (m *UnspentOutputsResponse) String() string            { return proto.CompactTextString(m) }
func (*UnspentOutputsResponse) ProtoMessage()               {}
func (*UnspentOutputsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *UnspentOutputsResponse) GetOutputs() []*WireUnspentOutput {
	if m != nil {
		return m.Outputs
	}
	return nil
}

type ProposeTransactionRequest struct {
	Transaction      *entities.Transaction `protobuf:"bytes,1,req,name=transaction" json:"transaction,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *ProposeTransactionRequest) Reset()                    { *m = ProposeTransactionRequest{} }
func (m *ProposeTransactionRequest) String() string            { return proto.CompactTextString(m) }
func (*ProposeTransactionRequest) ProtoMessage()               {}
func (*ProposeTransactionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ProposeTransactionRequest) GetTransaction() *entities.Transaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

type ProposeTransactionResponse struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *ProposeTransactionResponse) Reset()                    { *m = ProposeTransactionResponse{} }
func (m *ProposeTransactionResponse) String() string            { return proto.CompactTextString(m) }
func (*ProposeTransactionResponse) ProtoMessage()               {}
func (*ProposeTransactionResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*Account)(nil), "rpc.Account")
	proto.RegisterType((*WireUnspentOutput)(nil), "rpc.WireUnspentOutput")
	proto.RegisterType((*AccountRequest)(nil), "rpc.AccountRequest")
	proto.RegisterType((*AccountResponse)(nil), "rpc.AccountResponse")
	proto.RegisterType((*UnspentOutputsRequest)(nil), "rpc.UnspentOutputsRequest")
	proto.RegisterType((*UnspentOutputsResponse)(nil), "rpc.UnspentOutputsResponse")
	proto.RegisterType((*ProposeTransactionRequest)(nil), "rpc.ProposeTransactionRequest")
	proto.RegisterType((*ProposeTransactionResponse)(nil), "rpc.ProposeTransactionResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for WalletService service

type WalletServiceClient interface {
	Account(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountResponse, error)
	UnspentOutputs(ctx context.Context, in *UnspentOutputsRequest, opts ...grpc.CallOption) (*UnspentOutputsResponse, error)
	ProposeTransaction(ctx context.Context, in *ProposeTransactionRequest, opts ...grpc.CallOption) (*ProposeTransactionResponse, error)
}

type walletServiceClient struct {
	cc *grpc.ClientConn
}

func NewWalletServiceClient(cc *grpc.ClientConn) WalletServiceClient {
	return &walletServiceClient{cc}
}

func (c *walletServiceClient) Account(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	out := new(AccountResponse)
	err := grpc.Invoke(ctx, "/rpc.WalletService/Account", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) UnspentOutputs(ctx context.Context, in *UnspentOutputsRequest, opts ...grpc.CallOption) (*UnspentOutputsResponse, error) {
	out := new(UnspentOutputsResponse)
	err := grpc.Invoke(ctx, "/rpc.WalletService/UnspentOutputs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) ProposeTransaction(ctx context.Context, in *ProposeTransactionRequest, opts ...grpc.CallOption) (*ProposeTransactionResponse, error) {
	out := new(ProposeTransactionResponse)
	err := grpc.Invoke(ctx, "/rpc.WalletService/ProposeTransaction", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for WalletService service

type WalletServiceServer interface {
	Account(context.Context, *AccountRequest) (*AccountResponse, error)
	UnspentOutputs(context.Context, *UnspentOutputsRequest) (*UnspentOutputsResponse, error)
	ProposeTransaction(context.Context, *ProposeTransactionRequest) (*ProposeTransactionResponse, error)
}

func RegisterWalletServiceServer(s *grpc.Server, srv WalletServiceServer) {
	s.RegisterService(&_WalletService_serviceDesc, srv)
}

func _WalletService_Account_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Account(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.WalletService/Account",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Account(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_UnspentOutputs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnspentOutputsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).UnspentOutputs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.WalletService/UnspentOutputs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).UnspentOutputs(ctx, req.(*UnspentOutputsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_ProposeTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProposeTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).ProposeTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.WalletService/ProposeTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).ProposeTransaction(ctx, req.(*ProposeTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WalletService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.WalletService",
	HandlerType: (*WalletServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Account",
			Handler:    _WalletService_Account_Handler,
		},
		{
			MethodName: "UnspentOutputs",
			Handler:    _WalletService_UnspentOutputs_Handler,
		},
		{
			MethodName: "ProposeTransaction",
			Handler:    _WalletService_ProposeTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet.proto",
}

func init() { proto.RegisterFile("wallet.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 391 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x41, 0x4f, 0xdb, 0x30,
	0x18, 0x55, 0x93, 0xb5, 0x5d, 0xbf, 0x64, 0xdd, 0xe6, 0xad, 0x5d, 0x96, 0x55, 0x23, 0xca, 0x29,
	0xa7, 0x80, 0x22, 0x10, 0x42, 0x9c, 0xe0, 0x04, 0x5c, 0x40, 0xa6, 0x55, 0x25, 0x6e, 0x21, 0xb1,
	0x20, 0x52, 0x65, 0x07, 0xdb, 0x01, 0x0e, 0xfc, 0x5f, 0xfe, 0x06, 0x4a, 0xec, 0x34, 0x0d, 0x6d,
	0x11, 0xc7, 0x7c, 0xef, 0xf9, 0xbd, 0xf7, 0x3d, 0x3b, 0x60, 0x3f, 0xc5, 0x8b, 0x05, 0x91, 0x61,
	0xce, 0x99, 0x64, 0xc8, 0xe4, 0x79, 0xe2, 0xfe, 0x21, 0x34, 0x61, 0x69, 0x46, 0xef, 0x76, 0x09,
	0x95, 0x99, 0xcc, 0x88, 0x50, 0xa8, 0xff, 0x02, 0xfd, 0x93, 0x24, 0x61, 0x05, 0x95, 0xc8, 0x81,
	0x7e, 0x9c, 0xa6, 0x9c, 0x08, 0xe1, 0x74, 0x3c, 0x23, 0x18, 0xe0, 0xfa, 0x13, 0x4d, 0x60, 0x20,
	0x72, 0x42, 0xd3, 0xf8, 0x76, 0x41, 0x1c, 0xd3, 0x33, 0x82, 0x2f, 0xb8, 0x19, 0xa0, 0x23, 0xb0,
	0x25, 0x8f, 0xa9, 0x88, 0x13, 0x99, 0x31, 0x2a, 0x9c, 0xae, 0x67, 0x06, 0x56, 0x34, 0x0a, 0x97,
	0x4e, 0xd3, 0x06, 0xc5, 0x2d, 0xaa, 0xcf, 0xe0, 0xe7, 0x3c, 0xe3, 0x64, 0x46, 0x4b, 0x35, 0x79,
	0x59, 0xc8, 0xbc, 0xf8, 0x28, 0xc7, 0x18, 0x7a, 0xf2, 0xf9, 0x2c, 0x16, 0xf7, 0x8e, 0xe1, 0x19,
	0x81, 0x8d, 0xf5, 0x17, 0x0a, 0xa0, 0xc7, 0xaa, 0xb3, 0x55, 0x38, 0x2b, 0xfa, 0xd1, 0x78, 0x2b,
	0x4d, 0xac, 0x71, 0x3f, 0x84, 0xa1, 0x5e, 0x17, 0x93, 0x87, 0x82, 0x08, 0x59, 0xee, 0xa6, 0xe5,
	0x49, 0xe9, 0x67, 0x06, 0x03, 0xdc, 0x0c, 0xfc, 0x63, 0xf8, 0xbe, 0xe4, 0x8b, 0x9c, 0x51, 0x41,
	0x50, 0x00, 0x5f, 0x63, 0x35, 0x52, 0x7c, 0x2b, 0xb2, 0x43, 0x9e, 0x27, 0x61, 0xcd, 0x5b, 0xa2,
	0xfe, 0x01, 0x8c, 0x5a, 0x9b, 0x89, 0xcf, 0x79, 0x5e, 0xc0, 0xf8, 0xfd, 0x31, 0x6d, 0xbd, 0x07,
	0x7d, 0xb5, 0x47, 0xed, 0x3c, 0xae, 0x9c, 0xd7, 0x2a, 0xc4, 0x35, 0xcd, 0x9f, 0xc2, 0xdf, 0x2b,
	0xce, 0x72, 0x26, 0xc8, 0xea, 0x25, 0xe8, 0x18, 0x87, 0x60, 0xad, 0xdc, 0x46, 0x55, 0xf6, 0xd6,
	0x7b, 0x5b, 0x65, 0xfa, 0x13, 0x70, 0x37, 0xa9, 0xaa, 0x94, 0xd1, 0x6b, 0x07, 0xbe, 0xcd, 0xab,
	0x17, 0x78, 0x4d, 0xf8, 0x63, 0x96, 0x10, 0xb4, 0xdf, 0x3c, 0xb2, 0x5f, 0xad, 0xae, 0x54, 0x10,
	0xf7, 0x77, 0x7b, 0xa8, 0xb7, 0x3d, 0x87, 0x61, 0xbb, 0x07, 0xe4, 0x56, 0xbc, 0x8d, 0x9d, 0xba,
	0xff, 0x36, 0x62, 0x5a, 0x6a, 0x06, 0x68, 0x3d, 0x30, 0xfa, 0x5f, 0x1d, 0xd9, 0xda, 0x8f, 0xbb,
	0xb3, 0x15, 0x57, 0xb2, 0xa7, 0xdd, 0x9b, 0xf2, 0xe7, 0x7a, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x05,
	0x1a, 0x71, 0x48, 0x70, 0x03, 0x00, 0x00,
}
