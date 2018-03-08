package node

import (
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/blockrpc"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/closer"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/server/httpserver"
	"github.com/tclchiam/oxidize-go/server/rpc"
	walletRpc "github.com/tclchiam/oxidize-go/wallet/rpc"
)

type baseNode struct {
	p2p.PeerManager
	blockchain.Blockchain
	account.Engine

	rpcServer  *rpc.Server
	httpServer *httpserver.Server
}

func NewNode(bc blockchain.Blockchain, rpcServer *rpc.Server, httpServer *httpserver.Server) Node {
	return newNode(bc, rpcServer, httpServer)
}

func newNode(bc blockchain.Blockchain, rpcServer *rpc.Server, httpServer *httpserver.Server) *baseNode {
	node := &baseNode{
		Blockchain:  bc,
		PeerManager: p2p.NewPeerManager(),
		Engine:      account.NewEngine(bc),
		rpcServer:   rpcServer,
		httpServer:  httpServer,
	}

	blockrpc.RegisterSyncServer(rpcServer, blockrpc.NewSyncServer(bc))
	p2p.RegisterDiscoveryServer(rpcServer, p2p.NewDiscoveryServer(bc))
	walletRpc.RegisterWalletServer(rpcServer, walletRpc.NewWalletServer(node))

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

func (n *baseNode) Serve() {
	if n.rpcServer != nil {
		n.rpcServer.Serve()
	}
	if n.httpServer != nil {
		n.httpServer.Serve()
	}
}

func (n *baseNode) Close() error {
	closers := []closer.Closer{n.Blockchain, n.Engine}

	if n.rpcServer != nil {
		closers = append(closers, n.rpcServer)
	}
	if n.httpServer != nil {
		closers = append(closers, n.httpServer)
	}

	return closer.CloseMany(closers...)
}
