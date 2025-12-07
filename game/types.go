package game

import (
	"github.com/AndrewOPP/ZumaGameCoop/config"
	"github.com/AndrewOPP/ZumaGameCoop/hub"
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
	Hub   *hub.Hub
	Cfg   *config.Config
}

// MANAGER.GO
