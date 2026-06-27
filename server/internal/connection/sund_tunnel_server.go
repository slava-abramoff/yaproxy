package connection

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Отправка сообщений туннелю
func (h *ConnectionHandler) sendTunnelServer(clientID [4]byte, payload []byte) error {
	header := make([]byte, 8)

	copy(header[0:4], clientID[:])
	payloadLength := uint32(len(payload))
	binary.BigEndian.PutUint32(header[4:8], payloadLength)

	tunnelServer := h.hub.GetTunnel()
	if tunnelServer == nil {

		client, ok := h.hub.Get(clientID)
		if !ok {
			return fmt.Errorf("tunnelled resource is not connected. ClientID %s is not exits.", clientID)
		}

		return fmt.Errorf("tunnelled resource is not connected. Dropping client %s.", client.Conn.RemoteAddr())
	}

	buffers := net.Buffers{header, payload}
	_, err := buffers.WriteTo(tunnelServer)
	if err != nil {
		return err
	}

	return nil
}
