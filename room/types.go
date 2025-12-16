package room

import (
	// "log"
	"sync"
	"github.com/AndrewOPP/ZumaGameCoop/player"
)

// type PlayerCommand struct {
// 	PlayerID    string             `json:"player_id"`
// 	CommandType string             `json:"type"`
// 	Data        map[string]float64 `json:"data"`
// 	// Поле для передачи строковых данных, таких как цвет
// 	Payload map[string]string `json:"payload"`
// }

// type Player struct {
// 	ID string `json:"id"` 
// 	Nickname string `json:"nickname"` 
// 	Role string  `json:"role"` 
// 	RoomId string `json:"roomId"`
// 	Conn *websocket.Conn 
// 	Send chan []byte
// }



// type Hub struct {
// 	clients map[*websocket.Conn]bool

// 	Register chan *websocket.Conn

// 	Unregister chan *websocket.Conn

// 	StateUpdates chan []byte

// 	InputGate chan PlayerCommand
// }





type Room struct {
	ID string `json:"id"` 
	HostID string `json:"hostId"` 
	Players map[string] *player.Player `json:"players"` 
	StateUpdates chan []byte
	InputGate chan *player.PlayerCommand
	PlayersMutex sync.RWMutex
	RoomName string

}

func (r *Room) InputGateChan() chan *player.PlayerCommand {
    return r.InputGate 
}

// GetID реализует метод из интерфейса player.RoomContext
func (r *Room) GetID() string {
    return r.ID
}

type PlayerInfo struct {
    PlayerID       string `json:"playerid"`
    Nickname string `json:"nickname"`
    Role     string `json:"role"`
}