package room

import (
	"fmt"
	"encoding/json"
)


// for {
// 	select {
// 	case client := <-h.Register:
// 		h.clients[client] = true

// 	case client := <-h.Unregister:
// 		if _, ok := h.clients[client]; ok {
// 			delete(h.clients, client)
// 			client.Close()
// 		}

// 	case message := <-h.StateUpdates:
// 		for client := range h.clients {
// 			go client.WriteMessage(websocket.TextMessage, message)
// 		}

// 		// case command := <-h.InputGate:

// 		// 	log.Printf("Received command from player %s: %s", command.PlayerID, command.CommandType)

// 	}
// }

func (r *Room) GetPlayerInfoList() []map[string]string {
	r.PlayersMutex.RLock()
	defer r.PlayersMutex.RUnlock()

	playerList := make([]map[string]string, 0, len(r.Players))
	for _, player := range r.Players {
		playerList = append(playerList, map[string]string{
			"id": player.ID,
			"role": player.Role,
			"nickname": player.Nickname,
		})
	}

	return playerList
}


func (room *Room) BroadcastRoomUpdate() {
	update := struct {
        Type    string        `json:"type"`
        Players []map[string]string  `json:"players"` 
    }{
        Type:    "room_updated", // Тип, который ждет фронтенд
        Players: room.GetPlayerInfoList(), 
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

func (h *Room) Run() {

}

