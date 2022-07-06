package protocol

type StatsRequest struct {
	Protocol
}

type StatsReply struct {
	Protocol
	Height int    `json:"height"`
	Hash   string `json:"hash"`
}
