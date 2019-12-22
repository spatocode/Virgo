package core

import (
	"log"
	"time"
	"crypto/sha256"
	"encoding/hex"
)

var (
	genesisBlock = Block{0, "", "", time.Now().String(), ""}
	Blockchain = []Block{genesisBlock}
)

type Block struct {
	index		int
	hash		string
	prevHash	string
	timestamp	string
	data		string
} 

func CalculateHash(block Block) string {
	h := sha256.New()
	h.Write([]byte(string(block.index) + block.timestamp + block.prevHash + block.data))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func GenerateNextBlock(data string) Block {
	prevBlock := getLastBlock()
	nextIndex := prevBlock.index + 1
	nextHash := CalculateHash(prevBlock)
	nextTimestamp := time.Now().String()
	block := Block{
		index:		nextIndex,
		hash:		nextHash,
		prevHash:	prevBlock.hash,
		timestamp:	nextTimestamp,
		data:	data,
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
	if prevBlock.index + 1 != nextBlock.index {
		log.Println("Invalid index!")
		return false
	}
	
	if prevBlock.hash != nextBlock.prevHash {
		log.Println("Invalid previous hash!")
		return false
	} 
	
	if CalculateHash(nextBlock) != nextBlock.hash {
		log.Printf("Invalid hash: %s -> %s\n", CalculateHash(nextBlock), nextBlock.hash)
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
