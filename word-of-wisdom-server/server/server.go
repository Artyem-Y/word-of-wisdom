package server

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	maxConcurrentConnections = 100 // Maximum number of concurrent connections
	maxRequestsPerSecond     = 10  // Maximum number of requests per second from a single IP
	proofDifficulty          = 4   // Number of leading zeros for Proof of Work
	x                        = 1000000
)

var (
	connectionLimiter  *ConnectionLimiter
	clientTracker      *ClientTracker
	ddosMitigationChan chan net.IP
	quoteList          []string
)

// RunServer starts the Word of Wisdom Server with the provided quotes.
func RunServer(quotes []string) {
	quoteList = quotes

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)

		return
	}
	defer listener.Close()

	fmt.Println("Server is running. Waiting for connections...")

	connectionLimiter = NewConnectionLimiter(maxConcurrentConnections)
	clientTracker = NewClientTracker()

	ddosMitigationChan = make(chan net.IP, maxConcurrentConnections)

	go StartDDoSMitigation(ddosMitigationChan)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)

			continue
		}
		connectionLimiter.Acquire()
		go handleConnection(conn)
	}
}

// handleConnection handles an incoming client connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Get the client's IP address
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr).IP

	fmt.Println("remoteAddr", remoteAddr)

	if !clientTracker.CanProceed(remoteAddr) {
		ddosMitigationChan <- remoteAddr
		fmt.Println("Request rejected from", remoteAddr)

		return
	}

	// Limit on the number of concurrent connections
	connectionLimiter.Acquire()
	defer connectionLimiter.Release() // Ensure that the connection is released after processing

	// Check if the IP is blacklisted before proceeding
	if isIPBlacklisted(remoteAddr) {
		fmt.Println("IP address", remoteAddr, "is blacklisted, rejecting the connection.")

		return
	}

	// Generate challenge and send it to the client
	challenge := generateChallenge()
	challengeMessage := fmt.Sprintf("Solve Proof of Work: Find a nonce x such that SHA256(%d + %d) has %d leading zeros\n", x, challenge, proofDifficulty)
	_, err := conn.Write([]byte(challengeMessage))
	if err != nil {
		fmt.Println("conn.Write:", err)

		return
	}

	fmt.Println("challengeMessage", challengeMessage)

	// Read client nonce
	clientNonce, err := readClientNonce(conn)
	if err != nil {
		fmt.Println("Error reading client nonce:", err)

		return
	}
	fmt.Println("clientNonce", clientNonce)

	// Check Proof of Work
	if !checkProofOfWork(clientNonce) {
		fmt.Println("Connection rejected due to incorrect Proof of Work")

		return
	}

	// Send a confirmation and quote to the client
	confirmation := "Proof of Work verified. Connection approved."
	quote := getRandomQuote()
	send := fmt.Sprintf("%s:%s\n", confirmation, quote)
	_, err = conn.Write([]byte(send))
	if err != nil {
		fmt.Println("conn.Write:", err)

		return
	}

	fmt.Println("quote", quote)

	// Flush the connection to ensure the quote is sent before closing
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err = tcpConn.CloseWrite()
		if err != nil {
			fmt.Println("CloseWrite:", err)

			return
		}

		err = tcpConn.CloseRead()
		if err != nil {
			fmt.Println("CloseRead:", err)

			return
		}
	}
}

// readClientNonce reads the client nonce from the connection.
func readClientNonce(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	clientNonce, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(clientNonce), nil
}

// generateChallenge generates a random challenge value.
func generateChallenge() int {
	maxValue := 9000000
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(maxValue)
}

// checkProofOfWork checks if the client's Proof of Work is valid.
func checkProofOfWork(clientNonce string) bool {
	parts := strings.Split(clientNonce, ":")
	if len(parts) != 2 {
		return false
	}
	hashStr := parts[1]

	// Check that the hash string can be decoded
	_, err := hex.DecodeString(hashStr)
	if err != nil {
		fmt.Println("Error decoding hash:", err)

		return false
	}

	prefix := strings.Repeat("0", proofDifficulty)

	return strings.HasPrefix(hashStr, prefix)
}

// getRandomQuote returns a random quote from the quoteList.
func getRandomQuote() string {
	randomIndex := rand.Intn(len(quoteList))

	return quoteList[randomIndex]
}
