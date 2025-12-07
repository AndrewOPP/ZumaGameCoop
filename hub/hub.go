package hub

import (
	"github.com/gorilla/websocket"
)

// type PlayerCommand struct {
// 	PlayerID    string             `json:"player_id"`
// 	CommandType string             `json:"type"`
// 	Data        map[string]float64 `json:"data"`

// 	// Поле для передачи строковых данных, таких как цвет
// 	Payload map[string]string `json:"payload"`
// }

// type Hub struct {
// 	clients map[*websocket.Conn]bool

// 	Register chan *websocket.Conn

// 	Unregister chan *websocket.Conn

// 	StateUpdates chan []byte

// 	InputGate chan PlayerCommand
// }

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[*websocket.Conn]bool),
		Register:     make(chan *websocket.Conn),
		Unregister:   make(chan *websocket.Conn),
		InputGate:    make(chan PlayerCommand),
		StateUpdates: make(chan []byte),
	}
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}

		case message := <-h.StateUpdates:
			for client := range h.clients {
				go client.WriteMessage(websocket.TextMessage, message)
			}

			// case command := <-h.InputGate:

			// 	log.Printf("Received command from player %s: %s", command.PlayerID, command.CommandType)

		}
	}
}
