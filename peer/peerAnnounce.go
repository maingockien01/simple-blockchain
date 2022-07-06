package peer

import (
	"blockchain/protocol"
	"blockchain/udp"
	"encoding/json"
	"fmt"
)

// peer <- receive announce
func (peer *BlockchainPeer) HandleReciveAnnounce(udpMessage udp.UDPRequest) {
	announceByte := udpMessage.Message

	var announce protocol.AnnounceMessage

	json.Unmarshal(announceByte, &announce)

	var block Block

	block.Height = announce.Height
	block.Nonce = announce.Nonce
	block.Messages = announce.Messages
	block.Hash = announce.Hash
	block.MinedBy = announce.MinedBy

	peer.receiveAnnounceBlock(block)

}

func (peer *BlockchainPeer) receiveAnnounceBlock(block Block) {
	peer.isSyncing.Lock()
	isAdded := false
	if peer.Chain.VerifyNewBlock(block) {
		isAdded = peer.Chain.AddBlock(block)
	}
	fmt.Printf("Receive Announce:\n\tHash %s\n\tHeight %d\n\tisAdded: %t\n", block.Hash, block.Height, isAdded)

	peer.isSyncing.Unlock()
}

func (peer *BlockchainPeer) SendAnnounceRequest(block Block) {
	var announce protocol.AnnounceMessage
	announce.ProtocolType = "ANNOUNCE"

	announce.Hash = block.Hash
	announce.Height = block.Height
	announce.Messages = block.Messages
	announce.MinedBy = block.MinedBy
	announce.Nonce = block.Nonce

	jsonByte, _ := json.Marshal(announce)

	for _, contact := range peer.KnownPeers {
		go func(contact BlockchainPeerContact) {
			conn, _ := udp.OpenConnection(contact.Host, contact.Port)
			udp.SendMessage(conn, jsonByte, nil)
		}(contact)
	}
}
