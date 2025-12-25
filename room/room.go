package room

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/AndrewOPP/ZumaGameCoop/player"
)

func (room *Room) GetPlayerInfoList() []map[string]interface{} {
	room.PlayersMutex.RLock()
	defer room.PlayersMutex.RUnlock()

	playerList := make([]map[string]interface{}, 0, len(room.Players))

	for id, player := range room.Players {
		info := map[string]interface{}{
            "id":       player.ID,
            "nickname": player.Nickname,
            "role":     player.Role,
            "isReady":  room.State.ReadyStatus[id], 
            "score":    room.State.Scores[id],
        }

		playerList = append(playerList, info)
		// playerList = append(playerList, map[string]interface{}{
		// 	"id": player.ID,
		// 	"role": player.Role,
		// 	"nickname": player.Nickname,
		// })
	}

	return playerList
}

func (room *Room) BroadcastRoomUpdate() {
	room.Mu.Lock()
	state := room.State
	room.Mu.Unlock()

update := struct {
        Type            string                   `json:"type"`
        RoomID          string                   `json:"roomID"`
        RoomName        string                   `json:"roomName"`
        CurrentPlayerID string                   `json:"currentPlayerID"` // Оставь пустым или заполни, если нужно
        Players         []map[string]interface{} `json:"players"`
        GameState       GameState                `json:"gameState"` // Здесь будут маленькие буквы благодаря тегам в GameState
    }{
        Type:      "room_updated",
        RoomID:    room.ID,
        RoomName:  room.RoomName,
        Players:   room.GetPlayerInfoList(),
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

func (room *Room) makeAttempsWordsChars(playerID string) (string, map[string]string) {
	var attempsWords []string
	attempsWordsChars := make(map[string]string)	

	for _, attempt := range  room.State.PlayerAttempts[playerID] {
		attempsWords = append(attempsWords, attempt.Word)

		for index, char := range attempt.Word {
			attempsWordsChars[string(char)] = string(attempt.Result[index])
		}
	}

	allChars := strings.Join(attempsWords, "")
	return allChars, attempsWordsChars
}


func (room *Room) addScoreToPlayer(playerID string, charResult string, playersWordChar string, allCharsString string, allCharsMap map[string]string) {
	switch charResult{
	case  "G": 
		if !strings.Contains(allCharsString, playersWordChar){
			room.State.Scores[playerID]	+= 50
		} else if val, exists := allCharsMap[playersWordChar]; exists && val == "Y" {
			room.State.Scores[playerID]	+= 25
		}
	case "Y":
		if !strings.Contains(allCharsString, playersWordChar){
			room.State.Scores[playerID]	+= 10
		}
	}
}


func (room *Room) checkPlayersWord(playersWord string, secretWord string, playerID string) string {
	result := make([]string, 6) 
	usedInSecret := make([]bool, 6)
	allCharsString, allCharsMap := room.makeAttempsWordsChars(playerID)

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

return strings.Join(result, "")
}


func (room *Room) handleCommand(command *player.PlayerCommand) {
	room.Mu.Lock()
	playerID, _, commandType := command.PlayerID, command.Data, command.Type 
	

	switch commandType {
	case "toggle_ready": 
		currentStatus := room.State.ReadyStatus[playerID]
		room.State.ReadyStatus[playerID] = !currentStatus
		fmt.Printf("Игрок %s теперь %v", playerID, room.State.ReadyStatus[playerID])

		allReady := true 

		// 2. Проверяем всех. Если нашли хоть одного "не готового" — флаг падает
		for _, ready := range room.State.ReadyStatus {
			if !ready {
				allReady = false
				break 
			}
		}
		

		if allReady && !room.State.IsActive {
			room.State.IsActive = true
			fmt.Println("ВСЕ ИГРОКИ ГОТОВЫ, НАЧИНАЕМ ИГРУ")
			// Тут можно вызвать функцию старта: инициализировать слова, таймер и т.д.
			count := len(room.WordList)

			for _, player := range room.Players {
				randomIndex := rand.Intn(count)
				room.State.CurrentWords[player.ID] = room.WordList[randomIndex]
				fmt.Printf("Загаданное слово для игрока %s: %s\n ",   command.PlayerID,room.WordList[randomIndex])		
			}
		
    	} else {
			room.State.IsActive = false
		}

		
	case "check_word": 	
		var payload CheckWordPayload
		err := json.Unmarshal(command.Data, &payload)
		if err != nil {
			fmt.Printf("Error with  Unmarshal %s\n", err)
			// Если пришла абракадабра вместо строки — выходим или логируем
			return 
		}
		fmt.Printf("room.Dictionary[payload.Word] %s\n", room.Dictionary[payload.Word])
		if(room.Dictionary[payload.Word]) {

			fmt.Printf("YES WORD %s Is in the map \n", payload.Word)

			wordAnswer := room.checkPlayersWord(payload.Word, room.State.CurrentWords[playerID], playerID)	
			newAttempt := WordleAttempt{
				Word:    payload.Word,
				Result: wordAnswer, // Та самая строка "YYXXX"
			}
			room.State.PlayerAttempts[playerID] = append(room.State.PlayerAttempts[playerID], newAttempt)
			
		} else {
			fmt.Printf("WORD %s Is NOT in the map \n", payload.Word)
		}
		
	}
	room.Mu.Unlock()
	go room.BroadcastRoomUpdate()
}

