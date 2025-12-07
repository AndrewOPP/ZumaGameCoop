package constants

import (
	"fmt"
)

type CommandType string

const (
	// COMMAND_SHOOT - Отправить снаряд в указанную точку (x, y)
	CommandShoot CommandType = "SHOOT"

	// COMMAND_CHANGE_COLOR - Сменить цвет следующего снаряда (командная механика)
	CommandChangeColor CommandType = "CHANGE_COLOR"
)

func (c CommandType) String() string {
	switch c {
	case CommandShoot:
		return "Shoot"
	case CommandChangeColor:
		return "ChangeColor"
	default:
		// ВОТ ГДЕ ПРИХОДИТСЯ ИСПОЛЬЗОВАТЬ switch!
		return fmt.Sprintf("Unknown_Command(%d)", c)
	}
}
