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
	WebSocket WebSocketConfig
	// Physics PhysicsConfig
}

type WebSocketConfig struct {
	// Время ожидания Pong после отправки Ping (должно быть больше PingPeriod)
	PongWait time.Duration 

	// Период отправки Ping-сообщений клиенту (должен быть немного меньше PongWait)
	PingPeriod time.Duration 

	// Максимальный таймаут для операции записи (например, для WriteMessage или Ping)
	WriteWait time.Duration 

	// Максимальный размер входящего сообщения в байтах
	MaxMessageSize int64
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

		WebSocket: WebSocketConfig{
			MaxMessageSize: 512,            // 512 байт для команд JSON
            WriteWait:      10 * time.Second, // 10 секунд для таймаута записи
            PongWait:       60 * time.Second, // 60 секунд на ожидание Pong
            PingPeriod:     (60 * time.Second * 9) / 10, // 54 секунды
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
