package player

import "github.com/gorilla/websocket"

type PlayerCommand struct {
	PlayerID    string             `json:"player_id"`
	CommandType string             `json:"type"`
	Data        map[string]float64 `json:"data"`
	// Поле для передачи строковых данных, таких как цвет
	Payload map[string]string `json:"payload"`
}

type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	RoomId   string `json:"roomId"`
	Conn     *websocket.Conn
	Send     chan []byte
}
