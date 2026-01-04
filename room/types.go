package room

import (
	// "log"
	"sync"
	"github.com/AndrewOPP/ZumaGameCoop/player"
)


type Room struct {
	ID             string                      `json:"id"`
    HostID         string                      `json:"hostId"`
    RoomName       string                      `json:"roomName"`
    Players        map[string]*player.Player   `json:"players"`
    PlayersMutex   sync.RWMutex                // Защищает карту игроков
    
    InputGate      chan *player.PlayerCommand  // Канал для входящих команд
    StateUpdates   chan []byte                 // (Можно оставить для сырых данных)
    Done           chan struct{}               // Для остановки комнаты
    Dictionary     map[string]bool
    WordList      []string
    
    State          GameState                   // Состояние игры (очки, время)
    Mu             sync.Mutex                  // Защищает поле State

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

type WordleAttempt struct {
    Word   string `json:"word"`
    Result string `json:"result"` // Например, "GYXGG" (Green, Yellow, Gray)
    IsCorrect bool `json:"isCorrect"`
}
type Player struct {
    ID   string `json:"id"`
    Score int `json:"score"`
    Nickname string `json:"nickname"`
    Role     string `json:"role"`
    IsWaiting     bool `json:"isWaiting"`
    IsReady bool `json:"isReady"`

}

type GameState struct {
    CurrentWords   map[string]string          `json:"currentWords"`
    Players        map[string]*player.Player        `json:"players"`
    PlayerAttempts map[string][]WordleAttempt   `json:"playerAttempts"`
    TimeRemaining  int                        `json:"timeRemaining"`
    IsActive       bool                       `json:"isActive"`

}

type CheckWordPayload struct {
    Word string `json:"word"`
}