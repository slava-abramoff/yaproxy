package entities

import "net"

type Client struct {
	Conn net.Conn
	ID   [4]byte
}
