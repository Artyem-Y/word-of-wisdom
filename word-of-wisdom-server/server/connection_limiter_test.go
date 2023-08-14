package server

import (
	"sync"
	"testing"
	"time"
)

func TestConnectionLimiter_AcquireRelease(t *testing.T) {
	maxConnections := 3
	limiter := NewConnectionLimiter(maxConnections)

	// Attempt to acquire maxConnections slots concurrently
	var wg sync.WaitGroup
	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !limiter.Acquire() {
				t.Errorf("Acquire() returned false, expected true")
			}
		}()
	}

	wg.Wait()

	// Ensure that Acquire() returns false when maxConnections is reached
	if limiter.Acquire() {
		t.Errorf("Acquire() unexpectedly returned true after reaching maxConnections")
	}

	// Release one connection slot
	limiter.Release()

	// Acquire one more connection slot
	if !limiter.Acquire() {
		t.Errorf("Acquire() unexpectedly returned false after releasing a connection")
	}
}

func TestConnectionLimiter_ConcurrentAccess(t *testing.T) {
	maxConnections := 2
	limiter := NewConnectionLimiter(maxConnections)

	var wg sync.WaitGroup
	concurrency := 10
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Acquire() {
				defer limiter.Release()
				time.Sleep(time.Millisecond * 100) // Simulate some work
			}
		}()
	}

	wg.Wait()

	// Ensure that the connection count is correctly managed after concurrent access
	if limiter.connections != 0 {
		t.Errorf("Unexpected connection count, got: %d, expected: 0", limiter.connections)
	}
}
