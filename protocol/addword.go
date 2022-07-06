package protocol

type NewWordRequest struct {
	Protocol
	Word string `json:"word"`
}
