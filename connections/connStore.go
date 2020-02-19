package connections

import (
	"sync"
	"sync/atomic"
)

// ConnStore store all clients.
type ConnStore struct {
	// Connections concurrency safe map to save TCPConnection,
	// the key is int, the value is TCPConnection.
	connections sync.Map
	// length the count of all TCPConnection.
	length *int32
}

// Set store a TCPConnection.
func (cs *ConnStore) Set(conn *TCPConnection) {
	if _, ok := cs.connections.LoadOrStore(conn.ID, conn); !ok {
		return
	}

	atomic.AddInt32(cs.length, 1)
	cs.connections.Store(conn.ID, conn)
}

// Get get a TCPConnection by ID.
func (cs *ConnStore) Get(connectionID int) (*TCPConnection, bool) {
	connection, ok := cs.connections.Load(connectionID)
	if !ok {
		return nil, false
	}

	return connection.(*TCPConnection), true
}

// Del delete a TCPConnection by ID.
func (cs *ConnStore) Del(connectionID int) {
	cs.connections.Delete(connectionID)
	atomic.AddInt32(cs.length, -1)
}

// Range calls f sequentially for each key and value present in the map.
func (cs *ConnStore) Range(f func(key interface{}, value interface{}) bool) {
	cs.connections.Range(func(key, value interface{}) bool {
		return f(key, value)
	})
}

// Len return the count of all TCPConnection.
func (cs *ConnStore) Len() int32 {
	return *cs.length
}

// NewConnStore return ConnStore instance.
func NewConnStore() *ConnStore {
	return &ConnStore{
		connections: sync.Map{},
		length:      new(int32),
	}
}
