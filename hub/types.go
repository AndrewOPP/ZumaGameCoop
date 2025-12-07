package hub

import (
	"github.com/gorilla/websocket"
)

type PlayerCommand struct {
	PlayerID    string             `json:"player_id"`
	CommandType string             `json:"type"`
	Data        map[string]float64 `json:"data"`

	// Поле для передачи строковых данных, таких как цвет
	Payload map[string]string `json:"payload"`
}

type Hub struct {
	clients map[*websocket.Conn]bool

	Register chan *websocket.Conn

	Unregister chan *websocket.Conn

	StateUpdates chan []byte

	InputGate chan PlayerCommand
}
