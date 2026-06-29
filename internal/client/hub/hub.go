package hub

import (
	"net"
	"sync"
)

type TunnelClientHub struct {
	tunnel    net.Conn
	localConn map[[4]byte]net.Conn
	mu        sync.RWMutex
}

func NewTunnelClientHub() *TunnelClientHub {
	return &TunnelClientHub{
		localConn: make(map[[4]byte]net.Conn),
	}
}

// Добавляем новое локальное соединение
func (h *TunnelClientHub) AddLocalConn(clientID [4]byte, conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.localConn[clientID] = conn
}

// Удаляем локальное соединение
func (h *TunnelClientHub) RemoveLocalConn(clientID [4]byte, conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	oldConn, ok := h.localConn[clientID]

	if !ok {
		return
	}

	if conn != oldConn {
		return
	}

	delete(h.localConn, clientID)
}

// Получаем локальное соединеие по id
func (h *TunnelClientHub) GetLocalConn(clientID [4]byte) (net.Conn, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	c, exists := h.localConn[clientID]
	return c, exists
}

// Установить туннель
func (h *TunnelClientHub) SetTunnel(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.tunnel = conn
}

// Получить туннель
func (h *TunnelClientHub) GetTunnel() net.Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.tunnel
}

// Удалить туннель
func (h *TunnelClientHub) RemoveTunnel(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.tunnel != conn {
		return
	}

	h.tunnel = nil
}
