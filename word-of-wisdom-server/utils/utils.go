package utils

import (
	"bufio"
	"fmt"
	"os"
)

// ReadQuotesFromFile reads quotes from a text file and returns them as a slice of strings.
func ReadQuotesFromFile(filename string) []string {
	// Open the specified file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file:", err)

		return nil
	}
	defer file.Close()

	var quotes []string
	scanner := bufio.NewScanner(file)

	// Read each line from the file and add it to the quotes slice
	for scanner.Scan() {
		quotes = append(quotes, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading the file:", err)

		return nil
	}

	return quotes
}
