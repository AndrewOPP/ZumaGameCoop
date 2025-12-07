package game

import (
	"github.com/AndrewOPP/ZumaGameCoop/constants"
	"github.com/AndrewOPP/ZumaGameCoop/hub"
	"log"
)

func (gm *GameManager) HandleCommand(cmd hub.PlayerCommand) {
	commandType := constants.CommandType(cmd.CommandType)
	// Подсказка 3: Проверка совпадения command.Type
	if commandType == "CHANGE_COLOR" {
		newColor := cmd.Payload["color"]
		log.Println("newColor:", newColor)
		if gm.State.CurrentBall != nil {
			// 3. Корректное присваивание цвета
			gm.State.CurrentBall.Color = newColor
		}

	}
}
