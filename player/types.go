package player

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"sync"
)


// RoomContext определяет функциональность, которую Room должен
// предоставить игроку (и его ReadPump) для маршрутизации команд.
type RoomContext interface {
    // InputGate возвращает канал для отправки команд.
    InputGateChan() chan *PlayerCommand 
    
    // GetID возвращает ID комнаты для логирования.
    GetID() string

    // Добавь другие методы, если они нужны в ReadPump (например, Unregister)
}


type PlayerCommand struct {
	PlayerID    string             `json:"player_id"`
	Type string             `json:"type"`
	Data        json.RawMessage `json:"data"`
}

type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	RoomId   string `json:"roomId"`
    Score int `json:"score"`
    IsWaiting     bool `json:"isWaiting"`
    IsReady bool `json:"isReady"`
	Conn      *websocket.Conn `json:"-"` 
    Send      chan []byte     `json:"-"`
    Done      chan struct{}   `json:"-"`
	Mutex sync.Mutex 			`json:"-"`

}

// type Room struct {
// 	ID string `json:"id"` 
// 	HostID string `json:"hostId"` 
// 	Players map[string] *Player `json:"players"` 
// 	StateUpdates chan []byte
// 	InputGate chan *PlayerCommand

// }

type PlayerRawCommand struct {
	// Type - строка, определяющая команду (например, "move", "fire", "ready")
	Type string `json:"type"` 

	// Data - используется для хранения специфичных данных команды.
	// json.RawMessage позволяет отложить десериализацию этой части.
	// Room.Run() позже десериализует это поле в нужную структуру (e.g., CannonMoveData).
	Data json.RawMessage `json:"data"` 
}

