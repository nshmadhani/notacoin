package network

import (
	"fmt"
	"log"
	"net"

	"nshmadhani.com/notcoin/blockchain/util"
)

type Server struct {
	Listener *net.TCPListener
}

func (s Server) Start(handleConnection func(*net.TCPConn)) {
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn)
	}

}

func NewServer(portAddress int) (Server, string) {

	laddr := util.ToAAddr(fmt.Sprintf("localhost:%d", portAddress))
	ln, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is Running", laddr.String())

	server := Server{ln}
	return server, laddr.String()

}
