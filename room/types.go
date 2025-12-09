package room

import (
	"log"
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



type MainHub struct {
    // Хранит все активные комнаты. Ключ — это уникальный ID комнаты (например, "lobby_1", "match_34").
    Rooms map[string]*Room 
    
    // Мьютекс для безопасного доступа к карте Rooms из нескольких горутин (подключение, отключение, создание комнат).
    mu sync.RWMutex 
	Logger *log.Logger
} 

type Room struct {
	ID string `json:"id"` 
	HostID string `json:"hostId"` 
	Players map[string] *player.Player `json:"players"` 
	StateUpdates chan []byte
	InputGate chan *player.PlayerCommand

}

