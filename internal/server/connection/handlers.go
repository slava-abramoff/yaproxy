package connection

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

// Обработка сообщений от туннелируемого ресурса
func (h *ConnectionHandler) handleTunnelConnection(conn net.Conn) {
	defer conn.Close()

	err := h.tunnelHandshake(conn)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	h.hub.SetTunnel(conn)
	defer h.hub.RemoveTunnel(conn)

	log.Println("Tunnel is alive.")

	headerBuf := make([]byte, 8)
	for {
		_, err := io.ReadFull(conn, headerBuf)
		if err != nil {
			log.Println("Tunnel died or closed.")
			return
		}

		var clientID [4]byte
		copy(clientID[:], headerBuf[0:4])
		payloadLength := binary.BigEndian.Uint32(headerBuf[4:8])

		client, exists := h.hub.Get(clientID)
		if !exists {
			log.Printf("Client %x not found, discarding %d bytes...\n", clientID, payloadLength)

			// !!! Читаем эти байты в "черную дыру" (io.Discard) очищаем поток
			_, _ = io.CopyN(io.Discard, conn, int64(payloadLength))
			continue
		}

		_, err = io.CopyN(client.Conn, conn, int64(payloadLength))
		if err != nil {
			log.Printf("Failed to send data to public client %x: %v\n", clientID, err)

			// Если сокет клиента умер удаляем его из хаба
			h.hub.Remove(clientID)
			client.Conn.Close()
			continue
		}
	}
}

// Обработка сообщений с открытого порта
func (h *ConnectionHandler) handlePublicConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Client %s has connected!\n", conn.RemoteAddr())

	client, err := h.hub.Add(conn)
	if err != nil {
		log.Printf("Connection issues: %v\n", err)
		return
	}
	defer h.hub.Remove(client.ID)

	buffer := make([]byte, 32*1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Client %s has disconnected.\n", conn.RemoteAddr())

			return
		}

		// clientMessage := string(buffer[:n])
		// log.Printf("Client send: %s", clientMessage)
		// conn.Write([]byte("Сервер получил твое сообщение!\n"))\

		err = h.sendTunnelServer(client.ID, buffer[:n])
		if err != nil {
			log.Printf("Failed to route data to tunnel for client %s: %v\n", conn.RemoteAddr(), err)
			return
		}
	}

}
