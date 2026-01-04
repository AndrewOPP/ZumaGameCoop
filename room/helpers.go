package room

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/AndrewOPP/ZumaGameCoop/player"
)

func (room *Room) makeAttempsWordsChars(playerID string, ) (string, map[string]string) {
	var attempsWords []string
	attempsWordsChars := make(map[string]string)	

	for _, attempt := range  room.State.PlayerAttempts[playerID] {
		attempsWords = append(attempsWords, attempt.Word)



		for index, char := range attempt.Word {
			attempsWordsChars[string(char)] = string(attempt.Result[index])
		}
	}

	allChars := strings.Join(attempsWords, "")
	fmt.Printf("allChars %s\n", allChars)
	fmt.Printf("room.State.PlayerAttempts[playerID] LENTH %s\n", len(room.State.PlayerAttempts[playerID]))
	return allChars, attempsWordsChars
}

func (room *Room) addScoreToPlayer(playerID string, charResult string, playersWordChar string, allCharsString string, allCharsMap map[string]string) {
	playerScore := room.State.Players[playerID].Score

	switch charResult{
	case  "G": 
		if !strings.Contains(allCharsString, playersWordChar){
			playerScore	+= 50
		} else if val, exists := allCharsMap[playersWordChar]; exists && val == "Y" {
			playerScore	+= 25
		}
	case "Y":
		if !strings.Contains(allCharsString, playersWordChar){
			playerScore	+= 10
		}
	}

	room.State.Players[playerID].Score = playerScore
}

func (room *Room) addPlayerAttempt(playerID string, payloadWord string) string {
		wordAnswer, isCorrect:= room.checkPlayersWord(payloadWord, room.State.CurrentWords[playerID], playerID)	
		newAttempt := WordleAttempt{
			Word:    payloadWord,
			Result: wordAnswer,// Та самая строка "YYXXX"
			IsCorrect: isCorrect,
		}
		// room.State.PlayerAttempts[playerID][] = append(room.State.PlayerAttempts[playerID], newAttempt)
		room.State.PlayerAttempts[playerID] = append(room.State.PlayerAttempts[playerID], newAttempt)

		return newAttempt.Result
}

func (room *Room) chooseRandomWords(players map[string]*player.Player) {
	count := len(room.WordList)

	for _, player := range players {
		randomIndex := rand.Intn(count)
		room.State.CurrentWords[player.ID] = room.WordList[randomIndex]
	}
}

func (room *Room) newRoundForPlayer(playerID string) {
		count := len(room.WordList)
		randomIndex := rand.Intn(count)
		room.State.CurrentWords[playerID] = room.WordList[randomIndex]
		room.State.PlayerAttempts[playerID] = make([]WordleAttempt, 0)
}
