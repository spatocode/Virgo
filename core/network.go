package core

import (
	"os"
	"io"
	"log"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Message struct {
	data string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handlePostBlockchain(w http.ResponseWriter, r *http.Request) {
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock := GenerateNextBlock(m.data)

	if isValidBlock(Blockchain[len(Blockchain)-1], newBlock) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func handleGetPeers(w http.ResponseWriter, r *http.Request) {
	
}

func handleConnectPeers(w http.ResponseWriter, r *http.Request) {
	
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, statusCode int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}

func RunServer() {
	port := os.Getenv("PORT")
	mux := mux.NewRouter().StrictSlash(false)
	mux.HandleFunc("/api/blocks", handleGetBlockchain).Methods("GET")
	mux.HandleFunc("/api/blocks", handlePostBlockchain).Methods("POST")
	mux.HandleFunc("/api/peers", handleGetPeers).Methods("GET")
	mux.HandleFunc("/api/peers", handleConnectPeers).Methods("POST")
	
	server := http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}

	log.Println("Listening on ", port)
	server.ListenAndServe()
}
