package game

import (
	"encoding/json"
	"github.com/AndrewOPP/ZumaGameCoop/config"
	"github.com/AndrewOPP/ZumaGameCoop/hub"
	"log"
	"time"
)

// type GameState struct {
// 	TestCoordinate float64
// 	CurrentBall    *Ball
// }

// type GameManager struct {
// 	State GameState
// 	Hub   *hub.Hub
// 	Cfg   *config.Config
// }

func NewGameManager(h *hub.Hub, cfg *config.Config) *GameManager {
	return &GameManager{
		State: GameState{
			TestCoordinate: 0.0, // Начинаем с 0

			CurrentBall: &Ball{
				Color:  cfg.Game.BallColor,
				Radius: 20.0,
				X:      400.0,
				Y:      300.0,
			},
		},

		Hub: h,
		Cfg: cfg,
	}
}

func (gm *GameManager) Run() {
	ticker := time.NewTicker(gm.Cfg.Game.TickRate)
	defer ticker.Stop()
	for {
		select {

		// 1. Обработка входящего ВВОДА ИГРОКА (команд)
		case command := <-gm.Hub.InputGate:
			// Здесь происходит ВСЯ ЛОГИКА ИГРЫ, связанная с вводом:
			// - Расчет траектории выстрела
			// - Вставка шара в цепь
			// - Смена цвета патрона

			gm.HandleCommand(command)
			jsonMessage, err := json.Marshal(gm.State)
			if err == nil {
				gm.Hub.StateUpdates <- jsonMessage
			}
			// log.Printf("jsonMessage %s", jsonMessage)
			log.Printf("Received input command: Type=%s, PlayerID=%s", command.CommandType, command.PlayerID)

		// 2. Тик для обновления ФИЗИКИ
		case <-ticker.C:
			// 2.1. Обновление состояния: Движение цепи, снарядов, таймеры
			gm.State.TestCoordinate += 0.1
			// 2.2. Сериализация
			jsonMessage, err := json.Marshal(gm.State)

			// 2.3. Рассылка: Отправляем готовый JSON в Хаб
			if err == nil {
				gm.Hub.StateUpdates <- jsonMessage
			}
		}
	}
}
