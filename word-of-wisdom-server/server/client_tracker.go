package server

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ClientTracker keeps track of clients' request history.
type ClientTracker struct {
	clients map[string]*Client
	mu      sync.Mutex
}

// Client holds information about a client's requests.
type Client struct {
	lastRequestTime time.Time
	requests        int
}

// NewClientTracker creates a new instance of ClientTracker.
func NewClientTracker() *ClientTracker {
	return &ClientTracker{
		clients: make(map[string]*Client),
	}
}

// CanProceed checks whether a client's request can proceed based on rate limits.
func (tracker *ClientTracker) CanProceed(ip net.IP) bool {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	client, exists := tracker.clients[ip.String()]
	if !exists {
		tracker.clients[ip.String()] = &Client{
			lastRequestTime: time.Now(),
			requests:        1,
		}

		return true
	}

	if time.Since(client.lastRequestTime) > time.Second {
		client.requests = 0
		client.lastRequestTime = time.Now()
	}

	client.requests++

	// Allow requests if the number of requests is within the limit for the current IP
	if client.requests <= maxRequestsPerSecond {
		return true
	}

	fmt.Println("client.requests:", client.requests, "IP address:", ip.String(), "exceeded rate limit")

	return false
}
