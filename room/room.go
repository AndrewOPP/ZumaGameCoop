package room

import (
	// "fmt"
	// "time"

	// "github.com/gorilla/websocket"
)


// func NewHub() *Hub {
// 	return &Hub{
// 		clients:      make(map[*websocket.Conn]bool),
// 		Register:     make(chan *websocket.Conn),
// 		Unregister:   make(chan *websocket.Conn),
// 		InputGate:    make(chan PlayerCommand),
// 		StateUpdates: make(chan []byte),
// 	}
// }

func (h *Room) Run() {

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
}


// func(h *MainHub) CreateRoom(HostConnection *Player) (string, error) {
// 	// 1. Блокируем мьютекс на запись (для безопасного доступа к карте Rooms)
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	// 2. Создаем уникальный ID комнаты (например, с помощью библиотеки uuid или просто счетчика)
//     // Для простоты используем текущее время и имя хоста для псевдо-уникальности
// 	roomID := fmt.Sprintf("room_%s_%d", HostConnection.ID, time.Now().UnixNano())

// 	newRoom := &Room{
// 		ID: roomID,
// 		HostID: HostConnection.ID,
// 		Players: make(map[string]*Player),
// 		StateUpdates: make(chan []byte),
// 		InputGate:    make(chan *PlayerCommand),
// 	}


// 	newRoom.Players[HostConnection.ID] = HostConnection

// 	h.Rooms[roomID] = newRoom

// 	go newRoom.Run()

// 	h.Logger.Printf("Новая комната создана. ID: %s, Хост: %s", roomID, HostConnection.ID)

// 	return roomID, nil
// }