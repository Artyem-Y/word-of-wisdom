package server

import (
	"net"
	"testing"
	"time"
)

func TestDDos_StartDDoSMitigation(t *testing.T) {
	// Create a channel for testing
	ddosMitigationChan := make(chan net.IP)

	// Start the DDoS mitigation
	go StartDDoSMitigation(ddosMitigationChan)

	// Add an IP address to the mitigation channel
	testIP := net.ParseIP("192.168.1.1")
	ddosMitigationChan <- testIP

	// Wait for a short time to allow the DDoS mitigation to take effect
	time.Sleep(time.Second)

	// Check if the IP address is blacklisted
	if !isIPBlacklisted(testIP) {
		t.Errorf("Expected IP address %s to be blacklisted", testIP)
	}
}

func TestDDos_TestDeleteExpiredIPs(t *testing.T) {
	// Add an IP address to the blacklist with a short duration
	testIP := net.ParseIP("192.168.1.2")
	blockedTime := time.Now().Add(-2 * time.Minute)
	blacklist[testIP.String()] = blockedTime

	// Call the deleteExpiredIPs function
	deleteExpiredIPs()

	// Check if the IP address has been removed from the blacklist
	if isIPBlacklisted(testIP) {
		t.Errorf("Expected IP address %s to be removed from the blacklist", testIP)
	}
}

func TestDDos_TestIsIPBlacklisted(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")

	// Test case: IP is not blacklisted
	blacklistMutex.Lock()
	delete(blacklist, ip.String())
	blacklistMutex.Unlock()
	if isIPBlacklisted(ip) {
		t.Error("Expected IP to not be blacklisted")
	}

	// Test case: IP is blacklisted and still blocked
	blacklistMutex.Lock()
	blacklist[ip.String()] = time.Now().Add(time.Minute) // Add blocked IP to the blacklist
	blacklistMutex.Unlock()
	if !isIPBlacklisted(ip) {
		t.Error("Expected IP to be blacklisted")
	}

	// Test case: IP is blacklisted and unblocked
	blacklistMutex.Lock()
	blacklist[ip.String()] = time.Now().Add(-time.Minute) // Add blocked IP to the blacklist
	blacklistMutex.Unlock()
	if isIPBlacklisted(ip) {
		t.Error("Expected IP to not be blacklisted")
	}
}

func TestDDos_TestAddIPToBlacklist(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")

	// Add IP to the blacklist
	addIPToBlacklist(ip)

	// Ensure the IP is in the blacklist
	if !isIPBlacklisted(ip) {
		t.Error("Expected IP to be blacklisted")
	}

	// Update IP in the blacklist
	addIPToBlacklist(ip)

	// Ensure the IP is still in the blacklist
	if !isIPBlacklisted(ip) {
		t.Error("Expected IP to be blacklisted after updating")
	}
}
