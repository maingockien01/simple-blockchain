package peer

import (
	"blockchain/protocol"
	"blockchain/udp"
	"encoding/json"
	"fmt"
	"net"
)

// peer -> stats -> other peers
func (peer *BlockchainPeer) SendStatsRequest(contact BlockchainPeerContact, statReply *protocol.StatsReply) error {
	var statsRequest protocol.StatsRequest

	statsRequest.ProtocolType = "STATS"

	statsRequestByte, _ := json.Marshal(statsRequest)

	conn, err := udp.OpenConnection(contact.Host, contact.Port)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = udp.SendMessage(conn, statsRequestByte, peer.handleStatsReply(contact, statReply))

	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// This was a timeout
			err = udp.SendMessage(conn, statsRequestByte, peer.handleStatsReply(contact, statReply))
			if err != nil {
				return err
			}
		} else if err != nil {
			// This was an error, but not a timeout
		}
		return err
	}

	return nil
}

func (peer *BlockchainPeer) handleStatsReply(contact BlockchainPeerContact, statReply *protocol.StatsReply) func([]byte) {
	return func(reply []byte) {
		json.Unmarshal(reply, statReply)
	}
}

// peer <- stats <- other peers
func (peer *BlockchainPeer) HandleStatsRequest(udpMessage udp.UDPRequest) {
	stats := peer.getStats()

	statsByte, _ := json.Marshal(stats)

	udpMessage.WriteBack(statsByte)

}

func (peer *BlockchainPeer) getStats() (stats protocol.StatsReply) {
	peerStats := peer.Chain.GetStats()

	stats.ProtocolType = "STATS_REPLY"
	stats.Height = peerStats.LongestHeight
	stats.Hash = peerStats.Hash

	return stats
}
