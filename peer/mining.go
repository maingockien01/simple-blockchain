package peer

import (
	"sync"
)

type Miner struct {
	isChanging      sync.Mutex
	MessagePool     []string
	MinerName       string
	MaxMessagesSize int
}

func (miner *Miner) AddWord(word string) {
	if len(word) <= 20 && len(word) > 0 {
		miner.isChanging.Lock()
		miner.MessagePool = append(miner.MessagePool, word)
		miner.isChanging.Unlock()
	}
}

func (miner *Miner) GetMessages() []string {
	for len(miner.MessagePool) == 0 {
		//Wait till there is messages
	}
	miner.isChanging.Lock()
	messagesSize := min(len(miner.MessagePool), miner.MaxMessagesSize)
	messages := miner.MessagePool[:messagesSize]
	miner.MessagePool = miner.MessagePool[messagesSize:]
	miner.isChanging.Unlock()

	return messages

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
