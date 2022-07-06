package peer

import (
	. "blockchain/protocol"
	"blockchain/udp"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const PING_TIMEOUT_DURATION = 60 //sec

//----------------Sending Flooding----------------------
func (peer *BlockchainPeer) JoinNetwork() {
	//Add well known peer: silicon.cs.umanitoba.ca:8999
	//Create Flood Request
	floodRequest := NewFloodRequest(peer.Host, peer.Port, peer.Name)
	//Flooding
	peer.Flooding(floodRequest)
}

func (peer *BlockchainPeer) Flooding(floodRequest FloodRequest) (err error) {
	//Create json protocol
	jsonByte, err := json.Marshal(floodRequest)
	if err != nil {
		return err
	}

	for _, contact := range peer.KnownPeers {
		//Create UDP dialer
		conn, err := udp.OpenConnection(contact.Host, contact.Port)
		if err != nil {
			continue
		}

		//Send flood request to all host
		err = udp.SendMessage(conn, jsonByte, nil)
	}

	return err
}

func (peer *BlockchainPeer) HandleFloodReply(udpMessage udp.UDPRequest) {
	//Unmashal reply
	var floodReply FloodReply
	jsonByte := udpMessage.Message
	err := json.Unmarshal(jsonByte, &floodReply)

	if err != nil {
		return
	}

	//Add contact to known peers
	peer.AddOtherPeer(floodReply.Host, floodReply.Port, floodReply.Name)
}

func (peer *BlockchainPeer) AddOtherPeer(host string, port int, name string) {
	newContact := BlockchainPeerContact{
		Host:     host,
		Port:     strconv.Itoa(port),
		Name:     name,
		LastPing: int(time.Now().Unix()),
	}
	peer.isUsingPeerList.Lock()

	for i := 0; i < len(peer.KnownPeers); i++ {
		now := int(time.Now().Unix())
		knownPeer := peer.KnownPeers[i]
		if knownPeer.Host == host && knownPeer.Port == strconv.Itoa(port) {
			knownPeer.LastPing = now
			peer.isUsingPeerList.Unlock()

			return
		}

		if (now - knownPeer.LastPing) > PING_TIMEOUT_DURATION {
			defer peer.RemovePeer(knownPeer)
		}

	}

	peer.KnownPeers = append(peer.KnownPeers, newContact)
	peer.isUsingPeerList.Unlock()

}

func (peer *BlockchainPeer) RemovePeer(contact BlockchainPeerContact) {
	peer.isUsingPeerList.Lock()
	for i, knownContact := range peer.KnownPeers {
		if knownContact.Host == contact.Host && knownContact.Port == contact.Port {
			peer.KnownPeers = append(peer.KnownPeers[:i], peer.KnownPeers[i+1:]...)
		}
	}
	peer.isUsingPeerList.Unlock()

}

//-------------------Flood--------------------

func (peer *BlockchainPeer) HandleFloodRequest(udpMessage udp.UDPRequest) {
	var floodReqest FloodRequest
	jsonByte := udpMessage.Message
	json.Unmarshal(jsonByte, &floodReqest)

	//Check if we receive this before
	if peer.isFloodAlreadyRepeat(floodReqest) {

	} else {
		//Add the peer into contact list
		peer.AddOtherPeer(floodReqest.Host, floodReqest.Port, "")
		//Response to it
		peer.replyFloodRequest(floodReqest)
		//Forward to other peers
		ForwardRepeat(floodReqest, peer)
	}

}

func (peer *BlockchainPeer) replyFloodRequest(floodRequest FloodRequest) {
	conn, err := udp.OpenConnection(floodRequest.Host, strconv.Itoa(floodRequest.Port))
	if err != nil {

		return
	}

	floodReply := NewFloodReply(peer.Host, peer.Port, peer.Name)

	floodReplyJson, _ := json.Marshal(floodReply)

	udp.SendMessage(conn, floodReplyJson, nil)
}

//Send flood message from peer to contact
func ForwardRepeat(floodMessage FloodRequest, peer *BlockchainPeer) error {
	fmt.Printf("Forwarding flood id %s to %d peers\n", floodMessage.Id, len(peer.KnownPeers))
	for i := 0; i < len(peer.KnownPeers); i++ {
		contact := peer.KnownPeers[i]
		forwardRepeatUDP(floodMessage, contact.Host, contact.Port)
	}

	return nil
}

func forwardRepeatUDP(floodMessage FloodRequest, host, port string) error {
	conn, err := udp.OpenConnection(host, port)
	if err != nil {

		return err
	}

	floodRequestJson, _ := json.Marshal(floodMessage)

	err = udp.SendMessage(conn, floodRequestJson, nil)
	//fmt.Printf("Forward flood %s to %s\n", floodRequestJson, conn.RemoteAddr())
	return err
}

func (peer *BlockchainPeer) isFloodAlreadyRepeat(floodMessage FloodRequest) bool {
	peer.peerLock.Lock()
	defer peer.peerLock.Unlock()

	for _, id := range peer.FloodIdRecords {
		if id == floodMessage.Id {
			return true
		}
	}
	peer.FloodIdRecords = append(peer.FloodIdRecords, floodMessage.Id)
	return false
}
