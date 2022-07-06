package peer

import (
	"blockchain/protocol"
)

type Consensus struct {
	isDoingConsensus bool
	statsList        []Stats
	//Lock             sync.Mutex
}

type Stats struct {
	LongestHeight  int
	Hash           string
	AgreedContacts []BlockchainPeerContact
}

func NewConsensus() (consensus Consensus) {
	consensus.isDoingConsensus = false

	return
}

func (consensus *Consensus) AddStatsReply(contact BlockchainPeerContact, statsReply protocol.StatsReply) {
	for i := 0; i < len(consensus.statsList); i++ {
		stats := consensus.statsList[i]
		if stats.Hash == statsReply.Hash && stats.LongestHeight == statsReply.Height {
			consensus.statsList[i].AgreedContacts = append(consensus.statsList[i].AgreedContacts, contact)
			return
		}
	}

	newStats := Stats{
		Hash:          statsReply.Hash,
		LongestHeight: statsReply.Height,
	}

	newStats.AgreedContacts = append(newStats.AgreedContacts, contact)

	consensus.statsList = append(consensus.statsList, newStats)

}

func (consensus *Consensus) GetLongestStats() (stats Stats) {
	mostAgreed := 0
	longestHeight := 0
	for i := 0; i < len(consensus.statsList); i++ {
		if consensus.statsList[i].LongestHeight > longestHeight {
			stats = consensus.statsList[i]
			longestHeight = stats.LongestHeight
			mostAgreed = len(consensus.statsList[i].AgreedContacts)
		} else if consensus.statsList[i].LongestHeight == longestHeight {
			if len(consensus.statsList[i].AgreedContacts) > mostAgreed {
				stats = consensus.statsList[i]
				mostAgreed = len(consensus.statsList[i].AgreedContacts)
				longestHeight = stats.LongestHeight
			}
		}
		//fmt.Printf("Geting longest height stats: %d - longestHeight %d - most agreed %d\n after %d\n", i, mostAgreed, longestHeight, consensus.statsList[i].LongestHeight)
	}

	return
}
