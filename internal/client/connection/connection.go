package connection

type TCPConnection struct {
	hub      string
	destAddr string
	destPort string
	srcAddr  string
	srcPort  string
}

// Запустить туннель
func (c *TCPConnection) Start() {}

// Рукопожатие с сервером
func (c *TCPConnection) tunnelHandshake() {}

// Извлекает заголовок и полезную нагрузку из последовательности байт
func (c *TCPConnection) getPayload() {}

// Склеить в байты с хедером и полезной нагрузкой для отправки
func (c *TCPConnection) formatForSend() {}
