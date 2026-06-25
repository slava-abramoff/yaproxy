package client

import (
	"crypto/rand"
	"fmt"
	"net"
	"server/internal/entities"
	"sync"
)

type ClientHub struct {
	mu      sync.RWMutex
	clients map[[4]byte]*entities.Client
	tunnel  net.Conn
}

func NewClientHub() *ClientHub {
	return &ClientHub{
		clients: make(map[[4]byte]*entities.Client),
	}
}

// Add генерирует 4-байтный ID, упаковывает соединение в Client и добавляет в хаб
func (h *ClientHub) Add(conn net.Conn) (*entities.Client, error) {

	var id [4]byte
	_, err := rand.Read(id[:])
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации ID: %w", err)
	}

	newClient := &entities.Client{
		Conn: conn,
		ID:   id,
	}

	// Блокируем мапу на запись
	h.mu.Lock()
	h.clients[id] = newClient
	h.mu.Unlock()

	return newClient, nil
}

// Remove удаляет клиента из хаба по его соединению
func (h *ClientHub) Remove(id [4]byte) {
	// Блокируем мапу на запись
	h.mu.Lock()
	delete(h.clients, id)
	h.mu.Unlock()
}

// Get возвращает клиента, если он есть в мапе
func (h *ClientHub) Get(id [4]byte) (*entities.Client, bool) {
	// Блокируем только на чтение (быстрый доступ)
	h.mu.RLock()
	c, exists := h.clients[id]
	h.mu.RUnlock()
	return c, exists
}

// SetTunnel устанавливаем туннелируемое подключение
func (h *ClientHub) SetTunnel(conn net.Conn) {
	h.mu.Lock()
	h.tunnel = conn
	h.mu.Unlock()
}

// GetTunnel устанавливаем туннелируемое подключение
func (h *ClientHub) GetTunnel() net.Conn {
	return h.tunnel
}
