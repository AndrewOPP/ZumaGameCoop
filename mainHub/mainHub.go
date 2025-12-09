package mainhub

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	
	"github.com/AndrewOPP/ZumaGameCoop/room"
	"github.com/AndrewOPP/ZumaGameCoop/player"
	"github.com/gorilla/websocket"
)

func NewMainHub() *MainHub {
	return &MainHub{
		Logger: log.New(os.Stdout, "[MAINHUB] ", log.Ldate|log.Ltime|log.Lshortfile),
		Rooms:  make(map[string]*room.Room),
		PlayerRoomMap: make(map[string]string),
		UnregisterRoom: make(chan *room.Room),
	}
}


func(h *MainHub) CreateRoom(HostConnection *player.Player) (string, error) {
	// 1. Блокируем мьютекс на запись (для безопасного доступа к карте Rooms)
	h.mu.Lock()
	defer h.mu.Unlock()

	// 2. Создаем уникальный ID комнаты (например, с помощью библиотеки uuid или просто счетчика)
    // Для простоты используем текущее время и имя хоста для псевдо-уникальности
	roomID := fmt.Sprintf("room_%s_%d", HostConnection.ID, time.Now().UnixNano())

	newRoom := &room.Room{
		ID: roomID,
		HostID: HostConnection.ID,
		Players: make(map[string]*player.Player),
		StateUpdates: make(chan []byte),
		InputGate: make(chan *player.PlayerCommand),
	}

	newRoom.Players[HostConnection.ID] = HostConnection

	h.Rooms[roomID] = newRoom

	h.Logger.Printf("Новая комната создана. ID: %s, Хост: %s", roomID, HostConnection.ID)

	return roomID, nil
}

func (h *MainHub) RoutePlayer(conn *websocket.Conn, r * http.Request) {
	player := player.CreatePlayer(conn, r)

	action := r.URL.Query().Get("action")
	roomID := r.URL.Query().Get("roomId")

	h.Logger.Printf("Игрок %s подключился. Действие: %s, Комната ID: %s", player.ID, action, roomID)
	go player.WritePump()
	// go player.ReadPump(h) /// доделать

	var err error
	switch action {
		case "create":
			h.Logger.Printf("Игрок %s просит создать новую комнату.", player.ID)
			player.Role = "host"
			_, err = h.CreateRoom(player)
		case "join":
			h.Logger.Printf("Игрок %s просит присоединиться к комнате %s.", player.ID, roomID)
			player.Role = "guest"
			// err = h.JoinRoom(player, roomID)
		default:
			err = fmt.Errorf("неизвестное или отсутствующее действие: %s", action)
		}
	
	if err != nil {
		h.Logger.Printf("Ошибка маршрутизации игрока %s: %v. Закрытие соединения.", player.ID, err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		conn.Close()
    }	

	}


