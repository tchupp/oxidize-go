package rpc

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type ConnectionManager interface {
	OpenConnection(address string) (*grpc.ClientConn, error)
	GetConnection(address string) (*grpc.ClientConn, error)
	CloseConnection(address string) error
}

type connectionManager struct {
	cache map[string]*grpc.ClientConn
	lock  sync.RWMutex
}

func NewConnectionManager() ConnectionManager {
	return &connectionManager{
		cache: make(map[string]*grpc.ClientConn),
	}
}

func (c *connectionManager) OpenConnection(address string) (*grpc.ClientConn, error) {
	conn, err := startInsecureConnection(address)
	if err != nil {
		return nil, err
	}

	c.lock.Lock()
	c.cache[address] = conn
	c.lock.Unlock()

	return conn, nil
}

func (c *connectionManager) GetConnection(address string) (*grpc.ClientConn, error) {
	c.lock.RLock()
	conn := c.cache[address]
	c.lock.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("connection not open for address: '%s'", address)
	}
	return conn, nil
}

func (c *connectionManager) CloseConnection(address string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if conn := c.cache[address]; conn != nil {
		delete(c.cache, address)

		return conn.Close()
	}

	return nil
}

func startInsecureConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTimeout(500*time.Millisecond), grpc.WithInsecure())
}
