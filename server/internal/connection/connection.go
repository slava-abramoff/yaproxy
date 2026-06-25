package connection

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"server/internal/entities"
)

type HubProvider interface {
	SetTunnel(conn net.Conn)
	GetTunnel() net.Conn
	Add(conn net.Conn) (*entities.Client, error)
	Get(id [4]byte) (*entities.Client, bool)
	Remove(id [4]byte)
}

type ConnectionHandler struct {
	PublicListener net.Listener
	TunnelListener net.Listener
	hub            HubProvider
	token          [4]byte
}

// Начинаем подключать клиентов по публичному порту
func StartPublicListener() {}

// Ждем подключение туннелируемого хоста
func StartTunnelListener() {}

// Проверяем туннелируемого клиента на подлинность
func (h *ConnectionHandler) tunnelHandshake(conn net.Conn) error {
	return nil
}

// Обработка сообщений от туннелируемого ресурса
func (h *ConnectionHandler) handleTunnelConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("Tunnel is alive.")

	h.hub.SetTunnel(conn)
	defer h.hub.SetTunnel(nil)

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
		log.Printf("Connection issues: %s\n", conn.RemoteAddr())
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

func (h *ConnectionHandler) sendTunnelServer(clientID [4]byte, payload []byte) error {
	header := make([]byte, 8)

	copy(header[0:4], clientID[:])
	payloadLength := uint32(len(payload))
	binary.BigEndian.PutUint32(header[4:8], payloadLength)

	tunnelServer := h.hub.GetTunnel()
	if tunnelServer == nil {

		client, ok := h.hub.Get(clientID)
		if !ok {
			return fmt.Errorf("Tunnelled resource is not connected. ClientID %s is not exits.\n", clientID)
		}

		return fmt.Errorf("Tunnelled resource is not connected. Dropping client %s.\n", client.Conn.RemoteAddr())
	}

	// Отправляем 8 байт заголовка в сеть
	_, err := tunnelServer.Write(header)
	if err != nil {
		return err
	}

	// Отправляем полезную нагрузку
	_, err = tunnelServer.Write(payload)
	return err
}
