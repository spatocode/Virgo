package core

import (
	"log"
	"time"
	"crypto/sha256"
	"encoding/hex"
)

var (
	genesisBlock = Block{
		Index:		0,
		Hash:		"",
		PrevHash:	"",
		Timestamp:	time.Now().String(),
		Data:		"",
	}
	Blockchain = []Block{genesisBlock}
)

type Block struct {
	Index		int
	Hash		string
	PrevHash	string
	Timestamp	string
	Data		string
} 

func CalculateHash(block Block) string {
	h := sha256.New()
	h.Write([]byte(string(block.Index) + block.Timestamp + block.PrevHash + block.Data))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func GenerateNextBlock(data string) Block {
	prevBlock := getLastBlock()
	nextIndex := prevBlock.Index + 1
	nextHash := CalculateHash(prevBlock)
	nextTimestamp := time.Now().String()
	block := Block{
		Index:		nextIndex,
		Hash:		nextHash,
		PrevHash:	prevBlock.Hash,
		Timestamp:	nextTimestamp,
		Data:	data,
	}

	return block
}

func getLastBlock() Block {
	return Blockchain[len(Blockchain)-1]
}

func getBlockchain() []Block {
	return Blockchain
}

func isValidBlock(prevBlock, nextBlock Block) bool {
	if prevBlock.Index + 1 != nextBlock.Index {
		log.Println("Invalid index!")
		return false
	}
	
	if prevBlock.Hash != nextBlock.PrevHash {
		log.Println("Invalid previous hash!")
		return false
	} 
	
	if CalculateHash(nextBlock) != nextBlock.Hash {
		log.Printf("Invalid hash: %s -> %s\n", CalculateHash(nextBlock), nextBlock.Hash)
		return false
	}

	return true
}

func isValidGenesis(block Block) bool {
	return block == genesisBlock
}

func isValidChain(blockchain []Block) bool {
	if !isValidGenesis(blockchain[0]) {
		return false
	}

	for i := 1; i < len(blockchain); i++ {
		if !isValidBlock(blockchain[i-1], blockchain[i]) {
			return false
		}
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if isValidChain(newBlocks) && len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
