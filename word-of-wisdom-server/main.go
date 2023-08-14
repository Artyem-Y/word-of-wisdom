package main

import (
	"fmt"
	"test/word-of-wisdom/server"
	"test/word-of-wisdom/utils"
)

const (
	quotesFile = "quotes.txt"
)

func main() {
	quoteList := utils.ReadQuotesFromFile(quotesFile)
	if quoteList == nil {
		fmt.Println("Error reading quotes from file.")

		return
	}

	server.RunServer(quoteList)
}
