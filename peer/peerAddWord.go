package peer

import (
	"blockchain/protocol"
	"blockchain/udp"
	"encoding/json"
	"fmt"
)

func (peer *BlockchainPeer) HandleAddWordRequest(udpRequest udp.UDPRequest) {
	var newWordReq protocol.NewWordRequest

	json.Unmarshal(udpRequest.Message, &newWordReq)
	peer.PeerMiner.AddWord(newWordReq.Word)

	fmt.Printf("Receive new word: %s\n", newWordReq.Word)
}
