package server

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	banDuration = 5 * time.Minute // Duration for IP address blocking (e.g., 5 minutes).

	blacklist      = make(map[string]time.Time) // Map to store blocked IP addresses and their unblocking time.
	blacklistMutex sync.Mutex                   // Mutex for synchronized access to the blacklist.
)

// StartDDoSMitigation starts the DDoS mitigation process by receiving IPs through a channel.
func StartDDoSMitigation(ddosMitigationChan chan net.IP) {
	fmt.Println("banDuration: ", banDuration)
	fmt.Println("blacklist: ", blacklist)

	ticker := time.NewTicker(banDuration) // Создаем метроном, чтобы проверять черный список периодически.
	defer ticker.Stop()

	for {
		select {
		case ip := <-ddosMitigationChan:
			if isIPBlacklisted(ip) {
				fmt.Printf("IP address %s is already blacklisted.\n", ip)

				continue
			}

			// Apply DDoS mitigation mechanism
			// For example, add the IP address to the blacklist for a certain duration.
			addIPToBlacklist(ip)

			fmt.Printf("DDoS mitigation applied for IP address: %s\n", ip)

		case <-ticker.C:
			// Периодически проверяем черный список и удаляем IP-адреса, которые истекли.
			deleteExpiredIPs()
		}
	}
}

func deleteExpiredIPs() {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	currentTime := time.Now()

	for ip, unblockTime := range blacklist {
		if currentTime.After(unblockTime) {
			delete(blacklist, ip)
			fmt.Printf("IP address %s removed from the blacklist.\n", ip)
		}
	}
}

// isIPBlacklisted checks if an IP address is blacklisted.
func isIPBlacklisted(ip net.IP) bool {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	unblockTime, exists := blacklist[ip.String()]
	if !exists {
		return false // IP is not blacklisted
	}

	if time.Now().Before(unblockTime) {
		return true // IP is blacklisted and still blocked
	}

	// IP was blacklisted but has been unblocked
	delete(blacklist, ip.String())

	return false
}

// addIPToBlacklist adds an IP address to the blacklist for a specific duration.
func addIPToBlacklist(ip net.IP) {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	// Check if the IP is already in the blacklist
	if _, exists := blacklist[ip.String()]; !exists {
		// Add the IP address to the blacklist
		blacklist[ip.String()] = time.Now().Add(banDuration)
		fmt.Printf("IP address %s added to the blacklist.\n", ip)
	} else {
		// Update the unblocking time if the IP is already in the blacklist
		blacklist[ip.String()] = blacklist[ip.String()].Add(banDuration)
		fmt.Printf("IP address %s updated in the blacklist.\n", ip)
	}

	fmt.Println("banDuration2: ", banDuration)
	fmt.Println("blacklist2: ", blacklist)
}

func PrintBlacklist() {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	fmt.Println("Blacklisted IP addresses:")
	for ip, unblockTime := range blacklist {
		fmt.Printf("IP: %s, Unblock Time: %s\n", ip, unblockTime.String())
	}
}
