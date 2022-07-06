package peer

import (
	"blockchain/protocol"
	"blockchain/udp"
	"encoding/json"
	"net"
)

// peer -> get block request -> other peers

func (peer *BlockchainPeer) SendGetBlockRequest(height int, block *Block) {
	for i := 0; i < len(peer.KnownPeers); i++ {
		contact := peer.KnownPeers[i]
		peer.sendGetBlockRequest(height, contact.Host, contact.Port, block)
	}
}

func (peer *BlockchainPeer) sendGetBlockRequest(height int, toHost string, toPort string, block *Block) {
	getBlockRequest := protocol.GetBlockRequest{}
	getBlockRequest.ProtocolType = "GET_BLOCK"
	getBlockRequest.Height = height

	requestJsonByte, _ := json.Marshal(getBlockRequest)

	conn, err := udp.OpenConnection(toHost, toPort)

	if err != nil {
		return
	}

	err = udp.SendMessage(conn, requestJsonByte, peer.HandleGetBlockReply(block))

	if e, ok := err.(net.Error); ok && e.Timeout() {
		// This was a timeout
		err = udp.SendMessage(conn, requestJsonByte, peer.HandleGetBlockReply(block))
	} else if err != nil {
		// This was an error, but not a timeout
	}
}

func (peer *BlockchainPeer) HandleGetBlockReply(block *Block) func(reply []byte) {
	return func(reply []byte) {
		var getBlockReply protocol.GetBlockReply

		err := json.Unmarshal(reply, &getBlockReply)

		if err != nil {
			return
		}

		if getBlockReply.ProtocolType != "GET_BLOCK_REPLY" {
			return
		}

		if getBlockReply.Hash == "None" || getBlockReply.Messages == nil || getBlockReply.Nonce == "NONE" || getBlockReply.MinedBy == "None" || getBlockReply.Hash == "null" {
			return
		}

		block.Hash = getBlockReply.Hash
		block.Height = getBlockReply.Height
		block.Messages = getBlockReply.Messages
		block.MinedBy = getBlockReply.MinedBy
		block.Nonce = getBlockReply.Nonce
	}
}

// peer <- receive block request <- other peer
func (peer *BlockchainPeer) HandleGetBlockRequest(udpRequest udp.UDPRequest) {
	request := udpRequest.Message
	var getBlockRequest protocol.GetBlockRequest

	err := json.Unmarshal(request, &getBlockRequest)

	if err != nil {
		return
	}

	block := peer.Chain.GetBlock(getBlockRequest.Height)

	replyGetBlock(block, udpRequest)

}

func replyGetBlock(block Block, udpRequest udp.UDPRequest) {
	var getBlockReply protocol.GetBlockReply
	getBlockReply.ProtocolType = "GET_BLOCK_REPLY"
	getBlockReply.Hash = block.Hash
	getBlockReply.Height = block.Height
	getBlockReply.Messages = block.Messages
	getBlockReply.MinedBy = block.MinedBy
	getBlockReply.Nonce = block.Nonce

	jsonByte, _ := json.Marshal(getBlockReply)

	udpRequest.WriteBack(jsonByte)
}
