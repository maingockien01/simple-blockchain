package protocol

type GetBlockRequest struct {
	Protocol
	Height int `json:"height"`
}

type GetBlockReply struct {
	Protocol
	Height   int      `json:"height"`
	MinedBy  string   `json:"minedBy"`
	Nonce    string   `json:"nonce"`
	Messages []string `json:"messages"`
	Hash     string   `json:"hash"`
}
