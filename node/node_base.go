package node

import (
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/blockrpc"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/closer"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/rpc"
	walletRpc "github.com/tclchiam/oxidize-go/wallet/rpc"
)

type baseNode struct {
	p2p.PeerManager
	*rpc.Server
	blockchain.Blockchain
	account.Engine
}

func NewNode(bc blockchain.Blockchain, server *rpc.Server) Node {
	return newNode(bc, server)
}

func newNode(bc blockchain.Blockchain, server *rpc.Server) *baseNode {
	node := &baseNode{
		Blockchain:  bc,
		PeerManager: p2p.NewPeerManager(),
		Server:      server,
		Engine:      account.NewEngine(bc),
	}

	blockrpc.RegisterSyncServer(server, blockrpc.NewSyncServer(bc))
	p2p.RegisterDiscoveryServer(server, p2p.NewDiscoveryServer(bc))
	walletRpc.RegisterWalletServer(server, walletRpc.NewWalletServer(node))

	return node
}

func (n *baseNode) AddPeer(address string) (*p2p.Peer, error) {
	peer, err := n.PeerManager.AddPeer(address)
	if err != nil {
		return nil, err
	}

	go startSyncHeadersFlow(peer, n.PeerManager, n)

	return peer, nil
}

func (n *baseNode) SpendableOutputs(address *identity.Address) (*utxo.OutputSet, error) {
	return n.Engine.SpendableOutputs(address)
}

func (n *baseNode) Close() error {
	return closer.CloseMany(n.Blockchain, n.Engine, n.Server)
}
