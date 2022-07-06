package udp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

const WRITE_DEADLINE_MS = 100
const READ_DEADLINE_MS = 2000

func SendMessage(conn net.Conn, message []byte, replyHandler func([]byte)) (err error) {
	//fmt.Printf("Sending %s to", string(message))
	//fmt.Println(conn.RemoteAddr())

	buffer := make([]byte, DEFAULT_BUFFER_SIZE)
	fmt.Fprintf(conn, string(message))

	//Set Read Deadline
	conn.SetReadDeadline(time.Now().Add(READ_DEADLINE_MS * time.Microsecond))

	byteRead, err := bufio.NewReader(conn).Read(buffer)
	defer conn.Close()

	if err != nil {
		if err == io.EOF {
			//Do nothing cuz connection is closed on other side
			return nil
		}
		//fmt.Printf("Send Message Error: %s\n", err)
		return err
	}

	//fmt.Printf("UDP client: recived %s - from %s\n", buffer, conn.RemoteAddr())
	if replyHandler != nil {
		replyHandler(buffer[:byteRead])
	}

	return err
}

func OpenConnection(toHost string, port string) (conn net.Conn, err error) {
	address := toHost + ":" + port
	//fmt.Printf(">>> \tAddress: %s\n", address)
	conn, err = net.Dial("udp", address)

	if err != nil {
		fmt.Println(err)
		return
	}
	//Set Write Deadline
	conn.SetWriteDeadline(time.Now().Add(WRITE_DEADLINE_MS * time.Microsecond))
	return
}

func OpenConnectionAddr(address string) (conn net.Conn, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	//fmt.Printf(">>> \tAddress: %s\n", address)

	if err != nil {
		return nil, err
	}

	conn, err = net.DialUDP("udp4", nil, udpAddr)

	return conn, err
}
