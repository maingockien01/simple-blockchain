package peer

import (
	"blockchain/udp"
	"fmt"
	"net"
	"sync"
	"time"
)

const KNOWN_PEER_HOST = "130.179.28.37"
const KNOWN_PEER_PORT = 8999
const KNOWN_PEER_NAME = "Rob Peer"

type BlockchainPeerContact struct {
	Host     string
	Port     string
	Name     string
	LastPing int
}

type BlockchainPeer struct {
	BlockchainPeerContact

	FloodIdRecords []string

	KnownPeers []BlockchainPeerContact //addresses of other peers
	peerLock   sync.Mutex

	udpServer udp.UDPServer

	handlers map[string]func(udp.UDPRequest)

	Chain           Blockchain
	isSyncing       sync.Mutex
	isUsingPeerList sync.Mutex
	consensus       Consensus

	PeerMiner *Miner
}

func NewBlockchainPeer(host string, port string, name string) *BlockchainPeer {
	var newPeer BlockchainPeer
	addr, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("Host %s is not recognized!\n", host)
	}

	newPeer.Host = addr[0].To4().String()
	newPeer.Port = port
	newPeer.Name = name

	newPeer.handlers = make(map[string]func(udp.UDPRequest))

	newPeer.consensus = NewConsensus()

	newPeer.Chain.Difficulty = 8

	var miner Miner
	miner.MinerName = "ImKevin"
	miner.MaxMessagesSize = 20

	newPeer.PeerMiner = &miner
	newPeer.FloodIdRecords = []string{}

	return &newPeer
}

func (peer *BlockchainPeer) Run() {
	//Set up handler for each kind of protocol
	peer.handlers["FLOOD-REPLY"] = peer.HandleFloodReply
	peer.handlers["FLOOD"] = peer.HandleFloodRequest
	peer.handlers["STATS"] = peer.HandleStatsRequest
	peer.handlers["GET_BLOCK"] = peer.HandleGetBlockRequest
	peer.handlers["CONSENSUS"] = peer.HandleConsensusRequest
	peer.handlers["ANNOUNCE"] = peer.HandleReciveAnnounce
	peer.handlers["NEW_WORD"] = peer.HandleAddWordRequest

	//Set up UDP server
	peer.startServer()

	//Send flood to host
	peer.AddOtherPeer(KNOWN_PEER_HOST, KNOWN_PEER_PORT, KNOWN_PEER_NAME)
	peer.AddOtherPeer(KNOWN_PEER_HOST, KNOWN_PEER_PORT, KNOWN_PEER_NAME)
	fmt.Println(peer.KnownPeers)

	go peer.Ping()

	time.Sleep(10 * time.Second)
	go peer.Sync()

	go peer.Mining()
}

func (peer *BlockchainPeer) Ping() {
	for {
		fmt.Printf("Ping to %d peers\n", len(peer.KnownPeers))
		peer.JoinNetwork()
		time.Sleep(30 * time.Second)
	}
}

func (peer *BlockchainPeer) Sync() {
	for {
		peer.DoConsensus()
		time.Sleep(20 * time.Minute)
	}
}

func (peer *BlockchainPeer) Mining() {
	for {
		time.Sleep(5 * time.Minute)
		peer.Mine()
	}
}
