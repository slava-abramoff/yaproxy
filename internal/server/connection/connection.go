package connection

import (
	"fmt"
	"net"
	"yaproxy/internal/server/entities"
)

type HubProvider interface {
	RemoveTunnel(conn net.Conn)
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

func NewTCPHandler(hub HubProvider) (*ConnectionHandler, error) {
	handler := &ConnectionHandler{
		hub: hub,
	}

	token, err := handler.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure token: %w", err)
	}
	handler.SetToken(token)

	return handler, nil
}
