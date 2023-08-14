package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUtils_TestReadQuotesFromFile(t *testing.T) {
	// Create a temporary test file with some quotes
	tempFile, err := ioutil.TempFile("", "test_quotes_*.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	quotes := []string{
		"Quote 1",
		"Quote 2",
		"Quote 3",
	}

	// Write quotes to the temporary test file
	for _, quote := range quotes {
		_, err := tempFile.WriteString(quote + "\n")
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}
	tempFile.Close()

	// Test ReadQuotesFromFile
	readQuotes := ReadQuotesFromFile(tempFile.Name())

	// Check if the returned quotes match the expected quotes
	if len(readQuotes) != len(quotes) {
		t.Errorf("Expected %d quotes, but got %d", len(quotes), len(readQuotes))
	}

	for i, quote := range quotes {
		if readQuotes[i] != quote {
			t.Errorf("Expected quote %d to be '%s', but got '%s'", i, quote, readQuotes[i])
		}
	}
}
