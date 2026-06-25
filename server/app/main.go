package main

import (
	"log"
	"net"
)

func main() {
	publicListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	tunnelListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	defer publicListener.Close()
	defer tunnelListener.Close()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Client %s is ocnnected\n", conn.RemoteAddr())

	buffer := make([]byte, 1024)

}

// Обработка публичных соединений
func listenPublicPort(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Connection failed: %v", err)
			continue
		}

		// обрабатываем коннект в отдельной горутине

	}
}

func listenTunnel(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Connection failed: %v", err)
			continue
		}

		// обрабатываем коннект в отдельной горутине

	}
}
