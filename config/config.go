package config

import (
	"log"
	"time"
)

// ServerConfig содержит настройки сервера
type ServerConfig struct {
	Port            string
	Host            string
	AllowedOrigins  []string
	FrontMainFolder string
}

// PhysicsConfig содержит настройки геймплея
type PhysicsConfig struct {
	InitialBallRadius float64
	ChainSpeed        float64
}

type GameConfig struct {
	TickRate   time.Duration // Используется для game.GameManager
	MaxPlayers int
	BallColor  string
}

// Config объединяет все настройки
type Config struct {
	Server ServerConfig
	Game   GameConfig
	// Physics PhysicsConfig
}

func LoadConfig() *Config {
	// 1. Инициализация значениями по умолчанию
	cfg := &Config{
		Server: ServerConfig{
			Port:            ":8080",
			Host:            "localhost",
			AllowedOrigins:  []string{"*"},
			FrontMainFolder: "static",
		},

		Game: GameConfig{
			TickRate:   time.Second / 30,
			MaxPlayers: 2,
			BallColor:  "Black",
		},

		// Physics: PhysicsConfig{
		// 	InitialBallRadius: 20.0,
		// 	ChainSpeed:        0.5,
		// },
	}

	// 2. Логика чтения переменных окружения (опционально)
	// (например, if os.Getenv("PORT") != "" { cfg.Server.Port = os.Getenv("PORT") })

	log.Println("Configuration loaded successfully.")
	return cfg
}
