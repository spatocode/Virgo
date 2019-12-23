package core

import (
	"log"
	"time"
	"strings"
	"crypto/sha256"
	"encoding/hex"
)

var (
	genesisBlock = Block{Timestamp:	time.Now()} // Other fields will have nil values
	Blockchain = []Block{genesisBlock}
)

type Block struct {
	Index		int
	Hash		string
	PrevHash	string
	Timestamp	time.Time
	Data		string
	Difficulty	int
	Nonce		int
} 

func CalculateHash(block Block) string {
	h := sha256.New()
	h.Write([]byte(string(block.Index) + block.Timestamp.String() + block.PrevHash + block.Data + string(block.Difficulty) + string(block.Nonce)))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func GenerateNextBlock(data string) Block {
	prevBlock := getLastBlock()
	nextIndex := prevBlock.Index + 1
	nextHash, nonce := findBlockHash(prevBlock)
	nextTimestamp := time.Now()
	block := Block{
		Index:		nextIndex,
		PrevHash:	prevBlock.Hash,
		Timestamp:	nextTimestamp,
		Data:		data,	
	}

	block.Difficulty = getDifficulty(Blockchain)
	block.Hash = nextHash
	block.Nonce = nonce

	return block
}

func findBlockHash(block Block) (string, int) {
	var nonce int
	for {
		hash := CalculateHash(block)
		if isValidHash(hash, block.Difficulty) {
			return hash, nonce
		}
		nonce++
	}
}

func getLastBlock() Block {
	return Blockchain[len(Blockchain)-1]
}

func getBlockchain() []Block {
	return Blockchain
}

func getDifficulty(blockchain []Block) int {
	lastBlock := getLastBlock()
	if lastBlock.Index % difficultyAdjustmentInterval == 0 && lastBlock.Index != 0 {
		return getAdjustedDifficulty(lastBlock, blockchain)
	}
	return lastBlock.Difficulty
}

func getAdjustedDifficulty(lastBlock Block, blockchain []Block) int {
	prevBlockAdjustment := blockchain[len(Blockchain) - difficultyAdjustmentInterval]
	expectedTime := blockGenerationInterval * difficultyAdjustmentInterval
	timeTaken := lastBlock.Timestamp.Sub(prevBlockAdjustment.Timestamp)

	if int(timeTaken) < expectedTime / 2 {
		return prevBlockAdjustment.Difficulty + 1
	}
	
	if int(timeTaken) > expectedTime * 2 {
		return prevBlockAdjustment.Difficulty - 1
	}

	return prevBlockAdjustment.Difficulty
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

func isValidHash(hash string, difficulty int) bool {
	hashInBinary := hash //TODO: Convert the hash to a binary format
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashInBinary, prefix)
}

//TODO: Validate timestamp to mitigate the attack of manipulating difficulty
//func isValidTimestamp(prevBlock, nextBlock Block) bool {
//	
//}

func replaceChain(newBlocks []Block) {
	if isValidChain(newBlocks) && len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
