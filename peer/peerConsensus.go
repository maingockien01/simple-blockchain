package peer

import (
	"blockchain/protocol"
	"blockchain/udp"
	"encoding/json"
	"fmt"
	"sync"
)

//peer -> consensus -> other peers
func (peer *BlockchainPeer) SendConsensusRequest() {
	for _, contact := range peer.KnownPeers {
		go func(contact BlockchainPeerContact) {
			conn, _ := udp.OpenConnection(contact.Host, contact.Port)
			consensusRequest := protocol.ConsensusRequest{}
			consensusRequest.ProtocolType = "CONSENSUS"
			json, _ := json.Marshal(consensusRequest)
			udp.SendMessage(conn, json, nil)
		}(contact)
	}
}

//peer <- consensus <- other peers
func (peer *BlockchainPeer) HandleConsensusRequest(udpMessage udp.UDPRequest) {
	peer.DoConsensus()

}

func (peer *BlockchainPeer) DoConsensus() {
	//Start consensus
	peer.consensus = NewConsensus()
	peer.consensus.isDoingConsensus = true
	//Request stats from everyone
	var wg sync.WaitGroup
	peerToSend := len(peer.KnownPeers)
	//fmt.Printf("Sending to %s - %d\n", peer.KnownPeers, peerToSend)
	wg.Add(peerToSend)
	for i := 0; i < peerToSend; i++ {
		contact := peer.KnownPeers[i]
		go func(contact BlockchainPeerContact) {
			var statsReply protocol.StatsReply
			peer.SendStatsRequest(contact, &statsReply)
			//fmt.Printf("Do Consensus: get stats from %s:%s - height %d hash %s \n", contact.Host, contact.Port, statsReply.Height, statsReply.Hash)
			peer.consensus.AddStatsReply(contact, statsReply)
			defer wg.Done()
			//fmt.Println("DOne getting statsd")
		}(contact)
	}
	wg.Wait()
	// for _, stats := range peer.consensus.statsList {
	// 	fmt.Printf("After getting stats: hash %s - height %d - number of agreed contacts %d\n", stats.Hash, stats.LongestHeight, len(stats.AgreedContacts))
	// }
	//Choose the most agreed one
	stats := peer.consensus.GetLongestStats()
	fmt.Printf("Longest chain is %d height - %s hash\n", stats.LongestHeight, stats.Hash)

	peer.isSyncing.Lock()
	for peer.Chain.Height < stats.LongestHeight {
		blockHegihtNeedToGet := peer.Chain.Height
		isBlockAdded := false
		var block Block
		for _, contact := range stats.AgreedContacts {
			peer.sendGetBlockRequest(blockHegihtNeedToGet, contact.Host, contact.Port, &block)
			if peer.Chain.AddBlock(block) {
				isBlockAdded = true
				fmt.Printf("Added block: height %d - hash %s\n", block.Height, block.Hash)
				break
			}
		}
		if !isBlockAdded {
			peer.Chain.Height = 0 //Reset chain
		}

	}
	if !peer.Chain.VerifyBlockchain() {
		fmt.Println("Chain is not verified")
		peer.DoConsensus()
	}
	peer.isSyncing.Unlock()

	fmt.Println("Finish Do Consensus")
}
