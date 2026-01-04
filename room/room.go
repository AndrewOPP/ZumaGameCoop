package room

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/AndrewOPP/ZumaGameCoop/player"
)

// func (room *Room) GetPlayerInfoList() []map[string]interface{} {
// 	room.PlayersMutex.RLock()
// 	defer room.PlayersMutex.RUnlock()

// 	playerList := make([]map[string]interface{}, 0, len(room.Players))

// 	for id, player := range room.Players {
// 		info := map[string]interface{}{
//             "id":       player.ID,
//             "nickname": player.Nickname,
//             "role":     player.Role,
//             "isReady":  room.State.ReadyStatus[id], 
//             "score":    room.State.Scores[id],
//         }

// 		playerList = append(playerList, info)
// 	}

// 	return playerList
// }

func (room *Room) BroadcastRoomUpdate() {
	room.Mu.Lock()
	state := room.State
	room.Mu.Unlock()

update := struct {
        Type            string                   `json:"type"`
        RoomID          string                   `json:"roomID"`
        RoomName        string                   `json:"roomName"`
        CurrentPlayerID string                   `json:"currentPlayerID"` // Оставь пустым или заполни, если нужно
    
        GameState       GameState                `json:"gameState"` // Здесь будут маленькие буквы благодаря тегам в GameState
    }{
        Type:      "room_updated",
        RoomID:    room.ID,
        RoomName:  room.RoomName,
        GameState: state, // Вся логика (время, статус, попытки) теперь внутри этого объекта
    }

	message, err := json.Marshal(update)
	if err != nil {
		fmt.Printf("Ошибка маршалинга обновления комнаты %s: %v\n", room.ID, err)
		return  
	}

	room.PlayersMutex.RLock()
    defer room.PlayersMutex.RUnlock()

	for _, player := range room.Players {
		select {
		case player.Send <- message:

		default:
		}

	}
}

func (room *Room) Run() {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {

		select {
		case command := <- room.InputGate:
			room.handleCommand(command)
			fmt.Printf("Игра в комнату %s пришла команда %s",   room.ID, command)
		case <- ticker.C:
			room.Mu.Lock()
			if room.State.IsActive {
				room.State.TimeRemaining--

				if room.State.TimeRemaining <= 0 {
					room.State.IsActive = false
					fmt.Printf("Игра в комнате %s завершена по времени", room.ID)
				}
				room.Mu.Unlock()
				room.BroadcastRoomUpdate()
			} else {
				room.Mu.Unlock()
			}

		case <- room.Done:
			return
		}
	}
}

func (room *Room) checkPlayersWord(playersWord string, secretWord string, playerID string) (string, bool) {
	result := make([]string, 6) 
	usedInSecret := make([]bool, 6)
	allCharsString, allCharsMap := room.makeAttempsWordsChars(playerID)
    isCorrect := true

	for i := 0; i < 5; i++ {
		if playersWord[i] == secretWord[i] {
			result[i] = "G"
			usedInSecret[i] = true	
			room.addScoreToPlayer(playerID, result[i], string(playersWord[i]), allCharsString, allCharsMap)
    	}
	}	

	for i := 0; i < 5; i++ {
		if result[i] == "G" {
        	continue
    	}	
        isCorrect = false       
		result[i] = "X"

		for j := 0; j < 5; j++ {

			if !usedInSecret[j] {
				if playersWord[i] == secretWord[j] {
					result[i] = "Y"
					usedInSecret[j]  = true
					room.addScoreToPlayer(playerID, result[i], string(playersWord[i]), allCharsString, allCharsMap)
					break
				}
			}
		}
	}

  

	return strings.Join(result, ""), isCorrect
}


func (room *Room) handleCommand(command *player.PlayerCommand) {
    room.Mu.Lock()
    // Разблокируем мьютекс автоматически при выходе из функции
    defer room.Mu.Unlock()

    playerID := command.PlayerID
    commandType := command.Type 

    switch commandType {
    case "toggle_ready": 
        if player, ok := room.State.Players[playerID]; ok {
            player.IsReady = !player.IsReady
            fmt.Printf("Игрок %s теперь готов: %v\n", playerID, player.IsReady)
        }

        allReady := true 
        for _, player := range room.State.Players {
            if !player.IsReady {
                allReady = false
                break 
            }
        }

        if allReady && !room.State.IsActive {
            room.State.IsActive = true
            fmt.Println("ВСЕ ИГРОКИ ГОТОВЫ, НАЧИНАЕМ ИГРУ")
            room.chooseRandomWords(room.Players)
        } else if !allReady {
            // Если кто-то снял готовность, игра не активна (если еще не началась)
            if !room.State.IsActive {
                room.State.IsActive = false
            }
        }
        // Отправляем обновление статусов готовности
        go room.BroadcastRoomUpdate()


    case "leave_room":
        
        if player, ok := room.State.Players[playerID]; ok {
        player.Conn.Close() 
        }
        delete(room.State.Players, playerID)
        go room.BroadcastRoomUpdate()

    case "check_word":  
		// time.Sleep(500 * time.Millisecond)
        var payload CheckWordPayload
        if err := json.Unmarshal(command.Data, &payload); err != nil {
            fmt.Printf("Error with Unmarshal: %s\n", err)
            return 
        }

        // Проверяем, не находится ли игрок уже в режиме ожидания
        if room.State.Players[playerID].IsWaiting {
            return
        }
		
		
        if room.Dictionary[payload.Word] {
            fmt.Printf("Слово %s найдено в словаре\n", payload.Word)
            playerResult := room.addPlayerAttempt(playerID, payload.Word)
            
            if playerResult == "GGGGG" || len(room.State.PlayerAttempts[playerID]) >= 6 {
                // 1. Немедленно ставим флаг блокировки
               
                
                // 2. СРАЗУ рассылаем состояние, чтобы фронтенд заблокировал ввод и начал анимацию
                go room.BroadcastRoomUpdate()

                // 3. Запускаем таймер в отдельной горутине
                go func(playerID string) {
                    time.Sleep(4500*time.Millisecond)
                    
                    room.Mu.Lock()
                    // Проверяем, существует ли игрок (мог выйти за 3 секунды)
                    if player, ok := room.State.Players[playerID]; ok {
                        player.IsWaiting = false
                        room.newRoundForPlayer(playerID)
                    }
                    room.Mu.Unlock()
                    
                    // 4. Финальная рассылка после сброса раунда
                    room.BroadcastRoomUpdate() 
                }(playerID)
                
                // Выходим из кейса, так как Broadcast уже запущен в горутинеs
                return 
            }

        } else {
            fmt.Printf("Слова %s нет в словаре\n", payload.Word)
            // Здесь можно отправить игроку персональную ошибку "Not in word list"
        }
        
        // Рассылаем обновление (например, добавилась новая попытка, не победная)
        go room.BroadcastRoomUpdate()
    }
}



