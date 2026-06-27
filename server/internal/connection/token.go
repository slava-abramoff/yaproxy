package connection

import "server/utils/random"

// Генерируем новый токен
func (h *ConnectionHandler) GenerateToken() ([4]byte, error) {
	bytes, err := random.RandBytes(4)
	return [4]byte(bytes), err
}

// Под сомнением нахуй эта функция, но пусть будет, мало ли что там будет...
func (h *ConnectionHandler) SetToken(newToken [4]byte) {
	h.token = newToken
}

// Тоже не знаю зачем, но мало ли в будущем понадобится для бизнес логики
func (h *ConnectionHandler) GetToken() [4]byte {
	return h.token
}
