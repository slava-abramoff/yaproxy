package connection

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Рукопожатие туннеля
func (h *ConnectionHandler) tunnelHandshake(conn net.Conn) error {
	authBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, authBuf)
	if err != nil {
		return fmt.Errorf("handshake token not recognized: %v", err)
	}

	var incomingToken [4]byte
	copy(incomingToken[:], authBuf)

	if incomingToken != h.token {
		_, _ = conn.Write([]byte{0x00})
		return fmt.Errorf("access denied for %s, invalid token", conn.RemoteAddr())
	}

	_, err = conn.Write([]byte{0x01})
	if err != nil {
		return fmt.Errorf("failed to send handshake confirmation: %v", err)
	}

	log.Println("The handshake was successful; the tunnel is authorized!")

	return nil
}
