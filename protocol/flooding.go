package protocol

import (
	"strconv"

	"github.com/google/uuid"
)

type Flood struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
	Protocol
}

type FloodRequest struct {
	Id string `json:"id"`
	Flood
}

type FloodReply struct {
	Flood
}

func NewFloodReply(host string, port string, name string) (reply FloodReply) {
	reply.Host = host
	reply.Port, _ = strconv.Atoi(port)
	reply.Name = name
	reply.ProtocolType = "FLOOD-REPLY"

	return
}

func NewFloodRequest(host string, port string, name string) (request FloodRequest) {
	request.Host = host
	request.Port, _ = strconv.Atoi(port)
	request.Name = name
	request.ProtocolType = "FLOOD"
	request.Id = uuid.New().String()

	return request
}
