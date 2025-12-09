package game

import (
	"github.com/AndrewOPP/ZumaGameCoop/config"
	// "github.com/AndrewOPP/ZumaGameCoop/room"
)

type Ball struct {
	Color  string  `json:"color"`
	Radius float64 `json:"r"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

// MANAGER.GO

type GameState struct {
	TestCoordinate float64
	CurrentBall    *Ball
}

type GameManager struct {
	State GameState
	// Hub   *hub.Hub
	Cfg   *config.Config
}


// Player.GO



// Room.GO



// CreateRoomRequest - данные, которые клиент отправляет, чтобы создать комнату
type CreateRoomRequest struct {
	// Например, имя комнаты, которое выбрал пользователь
	RoomName string `json:"room_name"` 
	
	// ID или Никнейм игрока, который создал комнату и становится Хостом
	HostID string `json:"host_id"` 
	
	// Опционально: максимальное количество игроков
	MaxPlayers int `json:"max_players"` 
}

