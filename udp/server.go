package udp

import (
	"fmt"
	"net"
)

const DEFAULT_BUFFER_SIZE = 1024

type UDPServer struct {
	Port    string
	udpConn *net.UDPConn

	RequestHandler func(UDPRequest) //Must have
	ErrorHandler   func(error)      //Must have

	//Server config
}

type UDPRequest struct {
	Message    []byte
	RemoteAddr *net.UDPAddr
	UdpConn    *net.UDPConn
}

func (server *UDPServer) OpenSocket() error {
	fmt.Printf("Starting server with port %s\n", server.Port)
	udpAddr, err := net.ResolveUDPAddr("udp4", server.Port)
	fmt.Println(udpAddr)

	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		return err
	}

	server.udpConn = conn

	return nil
}

func (server *UDPServer) HandleSocket() {

	for {
		var buffer [DEFAULT_BUFFER_SIZE]byte

		n, addr, err := server.udpConn.ReadFromUDP(buffer[0:])

		if err != nil {
			server.ErrorHandler(err)
			continue
		}

		//fmt.Printf("UDP server received %s from %s\n", buffer[:n], addr)

		udpReq := UDPRequest{
			Message:    buffer[:n],
			RemoteAddr: addr,
			UdpConn:    server.udpConn,
		}

		go server.RequestHandler(udpReq)
	}
}

func (server *UDPServer) Stop() {
	server.udpConn.Close()
}

func (req UDPRequest) WriteBack(response []byte) {
	req.UdpConn.WriteToUDP(response, req.RemoteAddr)
}
