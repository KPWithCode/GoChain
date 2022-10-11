package main

import (
	"log"
	"net/http"
	"time"
)

const difficulty = 1

var GoBlock []Block

type Block struct {
	Index      int
	Timestamp  string
	Data       int
	Hash       string
	prevHash   string
	Difficulty int
	Nonce      string
}

func main() {
	err := godotdenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// goroutine for genesis block
	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}
	}()
	log.Fatal(run())
}

func run() error {

}

func makeRouter() http.Handler {
	HandleFunc("/", handleGetBlockchain).Methods("GET")
	HandleFunc("/", handleWriteBlock).Methods("POST")
}

func handleGetBlockchain() {

}

func handleWriteBlock() {

}

func respondwithJson() {

}

func generateBlock() {

}

func isBlockValid() bool {

}

func calculateHash() string {

}

func isHashValid() bool {

}
