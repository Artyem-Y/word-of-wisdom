package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"test/word-of-wisdom/word-of-wisdom-connection-check/conf"
)

const (
	proofDifficulty = 4 // Number of leading zeros for Proof of Work
)

func main() {
	if err := conf.Init(); err != nil {
		panic(err)
	}

	viper.SetDefault("settings.numConnections", 2) // Set the default value for "settings.numConnections"

	// Simulate multiple concurrent connections
	var numConnections = viper.GetInt("settings.numConnections")

	fmt.Println("numConnections:", numConnections)

	var wg sync.WaitGroup

	for i := 1; i <= numConnections; i++ {
		wg.Add(1)
		go simulateConnection(i, &wg)
	}

	wg.Wait()
}

func simulateConnection(connectionID int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("Connection %d error: %v\n", connectionID, err)

		return
	}
	defer conn.Close()

	challengeMessage, err := readChallenge(conn, connectionID)
	if err != nil {
		return
	}
	fmt.Printf("Connection %d received challenge: %s\n", connectionID, challengeMessage)

	challenge, err := extractAndSumChallengeValues(challengeMessage)
	if err != nil {
		fmt.Printf("Connection %d error extracting and summing challenge values: %v\n", connectionID, err)

		return
	}

	clientResponse := calculateProofOfWork(challenge)
	fmt.Printf("Connection %d calculated Proof of Work: %s\n", connectionID, clientResponse)

	_, err = conn.Write([]byte(clientResponse + "\n"))
	if err != nil {
		fmt.Printf("Connection %d error sending response: %v\n", connectionID, err)

		return
	}

	confirmationAndQuote, err := readConfirmationAndQuote(conn, connectionID)
	if err != nil {
		return
	}

	// Разделите подтверждение и цитату
	parts := strings.SplitN(confirmationAndQuote, ":", 2)
	if len(parts) != 2 {
		fmt.Printf("Connection %d error splitting confirmation and quote: invalid format\n", connectionID)

		return
	}
	confirmation := parts[0]
	quote := parts[1]
	fmt.Printf("Connection %d received confirmation: %s\n", connectionID, confirmation)
	fmt.Printf("Connection %d received quote: %s\n", connectionID, quote)
	fmt.Printf("---------------------------------------\n")
}

func readChallenge(conn net.Conn, connectionID int) (string, error) {
	reader := bufio.NewReader(conn)
	challengeMessage, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Connection %d error reading challenge: %v\n", connectionID, err)

		return "", err
	}

	return strings.TrimSpace(challengeMessage), nil
}

func extractAndSumChallengeValues(challengeMessage string) (int, error) {
	parts := strings.Split(challengeMessage, "+")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid challenge format")
	}
	valuesStr := strings.TrimSpace(parts[1])
	valuesStr = strings.TrimSuffix(valuesStr, ") has 4 leading zeros")

	values := strings.Fields(valuesStr)
	sum := 0
	for _, valStr := range values {
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return 0, err
		}
		sum += val
	}

	return sum, nil
}

func calculateProofOfWork(challenge int) string {
	leadingZeros := strings.Repeat("0", proofDifficulty)
	for nonce := 0; ; nonce++ {
		attempt := fmt.Sprintf("%d:%d", challenge, nonce)
		hash := sha256.Sum256([]byte(attempt))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, leadingZeros) {
			return fmt.Sprintf("%d:%s", nonce, hashStr)
		}
	}
}

func readConfirmationAndQuote(conn net.Conn, connectionID int) (string, error) {
	reader := bufio.NewReader(conn)
	confirmationAndQuote, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Connection %d error reading confirmation and quote: %v\n", connectionID, err)

		return "", err
	}

	return strings.TrimSpace(confirmationAndQuote), nil
}
