package protocol

type AnnounceMessage struct {
	Protocol
	Height   int      `json:"height"`
	MinedBy  string   `json:"minedBy"`
	Nonce    string   `json:"nonce"`
	Messages []string `json:"messages"`
	Hash     string   `json:"hash"`
}
