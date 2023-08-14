package server

import (
	"net"
	"testing"
)

func TestClientTracker_CanProceed(t *testing.T) {
	tracker := NewClientTracker()
	ip := net.ParseIP("127.0.0.1")

	// Test case: Successful requests
	for i := 0; i < maxRequestsPerSecond; i++ {
		if !tracker.CanProceed(ip) {
			t.Errorf("Expected CanProceed to return true for request %d", i)
		}
	}

	// Test case: Exceeded requests
	if tracker.CanProceed(ip) {
		t.Error("Expected CanProceed to return false after exceeding max requests")
	}
}
