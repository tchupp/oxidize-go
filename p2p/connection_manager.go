package p2p

import (
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ConnectionManager interface {
	OpenConnection(address string) (*grpc.ClientConn, error)
	HasConnection(address string) bool
	GetConnection(address string) *grpc.ClientConn
	CloseConnection(address string) error
}

type connectionEntry struct {
	conn  *grpc.ClientConn
	state connectivity.State
}

func (e *connectionEntry) WithState(state connectivity.State) *connectionEntry {
	e.state = state
	return e
}

type connectionManager struct {
	cache map[string]*connectionEntry
	lock  sync.RWMutex
}

func NewConnectionManager() ConnectionManager {
	return &connectionManager{
		cache: make(map[string]*connectionEntry),
	}
}

func (c *connectionManager) OpenConnection(address string) (*grpc.ClientConn, error) {
	conn, err := startInsecureConnection(address)
	if err != nil {
		return nil, err
	}

	c.saveEntry(address, &connectionEntry{conn: conn, state: conn.GetState()})
	go c.connectionMonitor(address, conn)

	return conn, nil
}

func (c *connectionManager) HasConnection(address string) bool {
	_, ok := c.cache[address]
	return ok
}

func (c *connectionManager) CloseConnection(address string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if entry := c.cache[address]; entry != nil {
		delete(c.cache, address)

		log.Infof("connection_manager: closing connection to '%s'", address)
		return entry.conn.Close()
	}

	return nil
}

func (c *connectionManager) saveEntry(address string, entry *connectionEntry) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[address] = entry
}

func (c *connectionManager) getEntry(address string) *connectionEntry {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache[address]
}

func (c *connectionManager) GetConnection(address string) *grpc.ClientConn {
	entry := c.getEntry(address)
	if entry == nil {
		return nil
	}
	return entry.conn
}

func (c *connectionManager) getConnectionState(address string) connectivity.State {
	entry := c.getEntry(address)
	if entry == nil {
		return -1
	}
	return entry.state
}

func (c *connectionManager) updateConnectionState(address string, newState connectivity.State) {
	entry := c.getEntry(address)
	if entry == nil {
		return
	}

	log.Debugf("connection_manager: connection state changed from '%s' to '%s'", entry.state, newState)
	c.saveEntry(address, entry.WithState(newState))
}

// Expected to be run in a goroutine
func (c *connectionManager) connectionMonitor(address string, conn *grpc.ClientConn) {
	for {
		currentState := conn.GetState()
		conn.WaitForStateChange(context.Background(), currentState)

		newState := conn.GetState()
		if newState == connectivity.TransientFailure {
			c.CloseConnection(address)
		}
		c.updateConnectionState(address, newState)
	}
}

func startInsecureConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address,
		grpc.WithTimeout(500*time.Millisecond),
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_logrus.StreamClientInterceptor(grpcLogger)),
		grpc.WithUnaryInterceptor(grpc_logrus.UnaryClientInterceptor(grpcLogger)),
	)
}
