package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const difficulty = 1

type Block struct {
	Index      int
	Timestamp  string
	Data       int
	Hash       string
	prevHash   string
	Difficulty int
	Nonce      string
}

var GoBlock []Block

type Message struct {
	Data int
}

var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// goroutine for genesis block
	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}
		spew.Dump(genesisBlock)
		mutex.Lock()
		GoBlock = append(GoBlock, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}

func run() error {
	mux := makeRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP server is running and listening on port:", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(GoBlock, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondwithJson(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	// sent old block and data
	newBlock := generateBlock(GoBlock[len(GoBlock) -1], m.Data)
	mutex.Unlock()
	// after block created check if block is valid
	if isBlockValid(newBlock, GoBlock[len(GoBlock) -1]) {
		GoBlock = append(GoBlock, newBlock)
		spew.Dump(GoBlock)
	}
	// json that shows new block created
	respondwithJson(w, r, http.StatusCreated, newBlock)
}

func respondwithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", " ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func generateBlock(oldBlock Block, Data int) Block {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.prevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("/%x", i)
		newBlock.Nonce = hex

		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), "You got more work to do!!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), "Hash work is complete")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	// new block has is update so return 
	return newBlock
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.prevHash != newBlock.prevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return false
}

func calculateHash(block Block) string {
	// data record
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Data) + block.prevHash + block.Nonce

	// hash it 
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	
	return hex.EncodeToString(hashed)
}

func isHashValid() bool {

}
