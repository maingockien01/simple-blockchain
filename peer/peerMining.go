package peer

import (
	"fmt"
	"strconv"
)

func (peer *BlockchainPeer) Mine() {
	//Get list of messages
	messages := peer.PeerMiner.GetMessages()
	//Do consensus
	peer.DoConsensus()

	//Prepare new block data
	lastBlock := peer.Chain.Chain[peer.Chain.Height-1]

	var newBlock Block
	newBlock.Height = peer.Chain.Height
	newBlock.Messages = messages
	newBlock.MinedBy = peer.PeerMiner.MinerName

	nonce := 1
	newBlock.Nonce = strconv.Itoa(nonce)
	newBlock.Hash = calculateHash(lastBlock, newBlock)
	//Finding nonce
	fmt.Println(">>> Start mining ...")
	for !isHashValid(newBlock.Hash, peer.Chain.Difficulty) {
		currentLastBlock := peer.Chain.Chain[peer.Chain.Height-1]

		if currentLastBlock.Height > lastBlock.Height {
			lastBlock = currentLastBlock
			nonce = 1
		}
		nonce++
		newBlock.Nonce = strconv.Itoa(nonce)
		newBlock.Hash = calculateHash(lastBlock, newBlock)
		//fmt.Printf("Try nonce %d\n - %s", nonce, newBlock.Hash)

	}
	//Announce
	fmt.Printf("Found nonce: %s - hash: %s\n", newBlock.Nonce, newBlock.Hash)
	peer.receiveAnnounceBlock(newBlock)
	peer.SendAnnounceRequest(newBlock)
}
