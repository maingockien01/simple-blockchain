package peer

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Blockchain struct {
	Chain      []Block
	Difficulty int
	Height     int
}

type Block struct {
	Height   int      `json:"height"`
	MinedBy  string   `json:"minedBy"`
	Nonce    string   `json:"nonce"`
	Messages []string `json:"messages"`
	Hash     string   `json:"hash"`
}

type InvalidNewBlockError struct {
	Message string
}

func (e *InvalidNewBlockError) Error() string {
	return e.Message
}

func (chain *Blockchain) AddBlock(block Block) bool {
	newBlockIndex := block.Height

	if newBlockIndex == 0 {
		chain.Chain = []Block{}
		chain.Chain = append(chain.Chain, block)
		chain.Height = 1
		return true
	}

	if block.Height > chain.Height || block.Height < 0 {
		return false
	}

	lastBlock := chain.Chain[newBlockIndex-1]

	if !VerifyBlock(lastBlock, block, chain.Difficulty) {
		return false
	}

	chain.Chain = chain.Chain[:lastBlock.Height+1]

	chain.Chain = append(chain.Chain, block)
	chain.Height = block.Height + 1

	return true
}

func GetNullBlock() Block {
	return Block{
		Height:   -1,
		MinedBy:  "NONE",
		Nonce:    "NONE",
		Messages: nil,
		Hash:     "NONE",
	}
}

func (chain *Blockchain) GetBlock(blockHeight int) Block {
	if blockHeight >= 0 && blockHeight < len(chain.Chain) {
		return chain.Chain[blockHeight]
	} else {
		return GetNullBlock()
	}
}

func (chain *Blockchain) VerifyBlockchain() bool {
	for i := 1; i < len(chain.Chain); i++ {
		lastBlock := chain.Chain[i-1]
		nextBlock := chain.Chain[i]
		if !VerifyBlock(lastBlock, nextBlock, chain.Difficulty) {
			return false
		}
	}
	return true
}

func VerifyBlock(lastBlock Block, nextBlock Block, difficult int) bool {
	if !isHashValid(nextBlock.Hash, difficult) {
		return false
	}
	if nextBlock.Height != lastBlock.Height+1 {
		return false
	}

	expectedHash := calculateHash(lastBlock, nextBlock)
	actualHash := nextBlock.Hash

	return expectedHash == actualHash
}

func isHashValid(hash string, difficulty int) bool {
	suffix := strings.Repeat("0", difficulty)
	return strings.HasSuffix(hash, suffix)
}

func calculateHash(lastBlock Block, block Block) string {
	blockHash := sha256.New()
	blockHash.Write([]byte(lastBlock.Hash))
	blockHash.Write([]byte(block.MinedBy))

	for _, message := range block.Messages {
		blockHash.Write([]byte(message))
	}
	blockHash.Write([]byte(block.Nonce))

	return fmt.Sprintf("%x", blockHash.Sum(nil))
}

func (chain *Blockchain) GetStats() Stats {
	longestHeight := chain.Height
	if longestHeight <= 0 {
		return Stats{
			LongestHeight: 0,
			Hash:          "NONE",
		}
	} else {
		return Stats{
			LongestHeight: chain.Height,
			Hash:          chain.Chain[chain.Height-1].Hash,
		}
	}
}

func (chain *Blockchain) VerifyNewBlock(block Block) bool {
	newBlockIndex := block.Height
	if block.Height != chain.Height {
		return false
	}
	if newBlockIndex == 0 {
		return true
	}

	lastBlock := chain.Chain[newBlockIndex-1]

	return VerifyBlock(lastBlock, block, chain.Difficulty)

}
