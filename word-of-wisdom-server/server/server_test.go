package server

import (
	"math/rand"
	"sort"
	"testing"
)

func TestGenerateChallenge(t *testing.T) {
	// Seed the random number generator to ensure consistent results
	rand.Seed(12345)

	// Generate a challenge
	challenge := generateChallenge()

	// Verify that the generated challenge is within the expected range
	minValue := 1000000
	maxValue := 1000000 + 9000000
	if challenge < minValue || challenge >= maxValue {
		t.Errorf("Generated challenge out of range, got: %d, expected range: [%d, %d)", challenge, minValue, maxValue)
	}
}

func TestCheckProofOfWork(t *testing.T) {
	clientNonce := "789:0000abcdef"
	if !checkProofOfWork(clientNonce) {
		t.Errorf("Proof of Work check failed for valid inputs")
	}

	invalidClientNonce := "789:abcdef0000"
	if checkProofOfWork(invalidClientNonce) {
		t.Errorf("Proof of Work check passed for invalid inputs")
	}
}

func TestGetRandomQuote(t *testing.T) {
	// Seed the random number generator to ensure consistent results
	rand.Seed(12345)

	// Create a list of test quotes
	testQuotes := []string{"Quote 1", "Quote 2", "Quote 3"}

	// Set the quoteList to the test quotes
	quoteList = testQuotes

	// Get a random quote
	quote := getRandomQuote()

	// Sort the test quotes and compare
	sort.Strings(testQuotes)
	if !contains(testQuotes, quote) {
		t.Errorf("Returned quote is not one of the test quotes, got: %s", quote)
	}
}

// Helper function to check if a slice contains a specific string.
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}
