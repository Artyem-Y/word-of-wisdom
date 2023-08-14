package server

import (
	"fmt"
	"sync"
	"time"
)

// ConnectionLimiter limits the number of concurrent connections.
type ConnectionLimiter struct {
	maxConnections int        // Maximum allowed concurrent connections.
	connections    int        // Current number of active connections.
	mu             sync.Mutex // Mutex to synchronize access to connection count.
}

// NewConnectionLimiter creates a new instance of ConnectionLimiter.
func NewConnectionLimiter(maxConnections int) *ConnectionLimiter {
	return &ConnectionLimiter{
		maxConnections: maxConnections,
		connections:    0,
	}
}

// Acquire attempts to acquire a connection slot and returns true if successful.
// If the maximum number of connections is reached, it waits and retries.
func (limiter *ConnectionLimiter) Acquire() bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// Check if the maximum number of connections has been reached
	for limiter.connections >= limiter.maxConnections {
		if limiter.connections >= limiter.maxConnections {
			fmt.Printf("Reached maximum connections (%d), waiting...\n", limiter.maxConnections)

			return false // Cannot acquire a connection slot at the moment
		}
		limiter.mu.Unlock()                // Unlock the mutex before sleeping
		time.Sleep(time.Millisecond * 100) // Pause to wait for a connection slot
		limiter.mu.Lock()                  // Lock the mutex again before retrying
	}

	// Increment the connection count and allow acquiring the slot
	limiter.connections++

	return true
}

// Release releases a connection slot.
func (limiter *ConnectionLimiter) Release() {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// Decrement the connection count to release the slot
	limiter.connections--
}
