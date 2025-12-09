package player

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func CreatePlayer(conn *websocket.Conn, r *http.Request) *Player  {
	playerID := uuid.New().String()
	nickname := r.URL.Query().Get("nickname")

	return &Player {
		ID: playerID,
		Nickname: nickname,
		Conn: conn,
		Send: make(chan []byte, 256),
		Role: "", 
        RoomId: "",
	}
}

func (player *Player) WritePump() {

	defer func() {
		player.Conn.Close()
	}()

	for message := range player.Send {
    // В Go for range автоматически выходит, когда канал закрыт,
    // и нет необходимости в проверке `if !ok`.

    err := player.Conn.WriteMessage(websocket.TextMessage, message)

    if err != nil {
        return // Возвращаемся при ошибке записи
    }
}
}

// func (player *Player) WritePump() {

// 	defer func() {
// 		player.Conn.Close()
// 	}()

// 	for {
// 	message, ok := <-player.Send

// 		if !ok {
// 			// player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 			return
// 		}

// 		err := player.Conn.WriteMessage(websocket.TextMessage, message)

// 		if err != nil {
// 			return
// 		}

// 	}
// }