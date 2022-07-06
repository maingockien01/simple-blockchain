package peer

import (
	"blockchain/protocol"
	. "blockchain/udp"
	"encoding/json"
	"fmt"
)

func (peer *BlockchainPeer) startServer() {
	peer.udpServer = UDPServer{
		Port:           ":" + peer.Port,
		RequestHandler: peer.protocolHandler,
		ErrorHandler:   peer.errorHandler,
	}

	peer.udpServer.OpenSocket()
	go peer.udpServer.HandleSocket()
}

func (peer *BlockchainPeer) protocolHandler(request UDPRequest) {
	//Parse protocol
	var protocol protocol.Protocol
	json.Unmarshal(request.Message, &protocol)

	//fmt.Printf("Recived message %s protocol type %s\n", string(request.Message), protocol.ProtocolType)

	if protocol.ProtocolType == "" {
		return
	}
	//Select further appropriate protocol handler based on protocol type
	handler := peer.handlers[protocol.ProtocolType]

	if handler == nil {
		//Invalid type
		return
	}

	handler(request)

}

func (peer *BlockchainPeer) errorHandler(err error) {
	//TODO:
	fmt.Printf("Error: %s\n", err)
}
