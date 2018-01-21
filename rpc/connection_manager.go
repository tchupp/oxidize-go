package rpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

type ConnectionManager struct {
	cache map[string]*grpc.ClientConn
	lock  sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		cache: make(map[string]*grpc.ClientConn),
	}
}

func (c *ConnectionManager) OpenConnection(address string) (*grpc.ClientConn, error) {
	conn, err := startInsecureConnection(address)
	if err != nil {
		return nil, err
	}

	c.lock.Lock()
	c.cache[address] = conn
	c.lock.Unlock()

	return conn, nil
}

func (c *ConnectionManager) CloseConnection(address string) error {
	c.lock.Lock()
	conn := c.cache[address]
	c.cache[address] = nil
	c.lock.Unlock()

	return conn.Close()
}

func startInsecureConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTimeout(500*time.Millisecond), grpc.WithInsecure())
}
