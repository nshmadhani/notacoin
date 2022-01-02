package network

import (
	"net"

	"nshmadhani.com/notcoin/blockchain/util"
)

type Client struct {
	Connection net.TCPConn
	RemoteAddr string
}

func (c *Client) Send(message []byte) {

	c.Connection.Write(message)
	c.Connection.Close()

}

func NewClient(addres string) *Client {
	conn, _ := net.DialTCP("tcp", nil, util.ToAAddr(addres))
	// if err != nil {
	// 	//log.Panic(err)
	// }

	return &Client{*conn, addres}
}
