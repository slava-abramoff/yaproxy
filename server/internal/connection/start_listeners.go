package connection

import (
	"log"
	"net"
	"strconv"
)

func (h *ConnectionHandler) newListener(port uint16) (net.Listener, error) {
	addr := ":" + strconv.Itoa(int(port))
	return net.Listen("tcp", addr)
}

// Начинаем подключать клиентов по публичному порту
func (h *ConnectionHandler) StartPublicListener(port uint16) {
	listener, err := h.newListener(port)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error: %v", err)
			conn.Close()
			continue
		}

		go h.handlePublicConnection(conn)
	}
}

// Ждем подключение туннелируемого хоста
func (h *ConnectionHandler) StartTunnelListener(port uint16) {
	listener, err := h.newListener(port)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error: %v", err)
			conn.Close()
			continue
		}

		go h.handleTunnelConnection(conn)
	}

}
