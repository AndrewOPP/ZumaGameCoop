package mainhub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/AndrewOPP/ZumaGameCoop/config"
	"github.com/AndrewOPP/ZumaGameCoop/player"
	"github.com/AndrewOPP/ZumaGameCoop/room"
	"github.com/gorilla/websocket"
)

func NewMainHub(cfg *config.Config) *MainHub {
	return &MainHub{
		Logger: log.New(os.Stdout, "[MAINHUB] ", log.Ldate|log.Ltime|log.Lshortfile),
		Rooms:  make(map[string]*room.Room),
		PlayerRoomMap: make(map[string]string),
		UnregisterRoom: make(chan *room.Room),
		Config: cfg,
	}
}

func(h *MainHub) CreateRoom(HostConnection *player.Player, roomName string) (*room.Room, error) {
	// 1. Блокируем мьютекс на запись (для безопасного доступа к карте Rooms)
	h.mu.Lock()
	defer h.mu.Unlock()

	// 2. Создаем уникальный ID комнаты (например, с помощью библиотеки uuid или просто счетчика)
    // Для простоты используем текущее время и имя хоста для псевдо-уникальности
	roomID := fmt.Sprintf("room_%s_%d", HostConnection.ID, time.Now().UnixNano())

	newRoom := &room.Room{
		ID: roomID,
		RoomName: roomName,
		HostID: HostConnection.ID,
		Players: make(map[string]*player.Player),
		StateUpdates: make(chan []byte),
		InputGate: make(chan *player.PlayerCommand),
	}

	newRoom.Players[HostConnection.ID] = HostConnection

	h.Rooms[roomID] = newRoom

	h.Logger.Printf("Новая комната создана. ID: %s, Хост: %s", roomID, HostConnection.ID)

	return newRoom, nil
}



func (h *MainHub) RoutePlayer(conn *websocket.Conn, r *http.Request) {
    // Получаем доступ к конфигу через поле структуры MainHub (предполагаем, что h.Config существует)
    cfg := h.Config // <-- ИСПРАВЛЕНИЕ: Получаем Config из Hub (если он там есть)
    // Если h.Config не существует, нужно добавить его в структуру MainHub.

	// player := player.CreatePlayer(conn, r)
	
	playerID := r.URL.Query().Get("playerID")
	action := r.URL.Query().Get("action")
	roomID := r.URL.Query().Get("roomId")
	roomName := r.URL.Query().Get("roomName")


    // 1. Объявляем room перед switch, чтобы она была доступна в return/if-блоках.
	var currentRoom *room.Room // Предполагаем, что *room.Room — это твой тип.
	var err error
	var messageType string
	var currentPlayer *player.Player

	h.Logger.Printf("Входящее подключение. Действие: %s, Комната ID: %s, Игрок ID: %s", action, roomID, playerID)

	switch action {
	case "create":
		currentPlayer = player.CreatePlayer(conn, r)
		h.Logger.Printf("Игрок %s просит создать новую комнату.", currentPlayer.ID)
		currentPlayer.Role = "host"
		messageType = "room_created"
        // 2. Присваиваем значение объявленной переменной room, используя '='
		currentRoom, err = h.CreateRoom(currentPlayer, roomName) 
	case "join":
		currentPlayer = player.CreatePlayer(conn, r)
		h.Logger.Printf("Игрок %s просит присоединиться к комнате %s.", currentPlayer.ID, roomID)
		currentPlayer.Role = "guest"
		messageType = "room_joined"
		currentRoom, err = h.JoinRoom(currentPlayer, roomID) // <-- Не забудь реализовать JoinRoom
		if(err != nil) {
			err = fmt.Errorf("неудалось подключиться, ошибка в подключении, %s", err)
		}
	case "reconnect":
		if(playerID == "" || roomID == "") {
			err = fmt.Errorf("отсутствуют необходимые ID для переподключения")
			break
		}
		h.Logger.Printf("Игрок %s просит переподключение к комнате %s.", playerID, roomID)

			currentPlayer, err = h.ReconnectPlayer(conn, roomID, playerID)

		// Если переподключение удалось, нам нужно также получить ссылку на комнату
		if err == nil && currentPlayer != nil {
			messageType = "room_reconnected"	
	
			// Предполагаем, что *Player знает свою комнату, или RoomHub может ее вернуть
			// Для простоты, здесь мы просто используем ID из запроса, если ReconnectPlayer не вернул ошибку:
			h.mu.RLock()
			currentRoom = h.Rooms[roomID] 
			h.mu.RUnlock()
		}

	default:
		err = fmt.Errorf("неизвестное или отсутствующее действие: %s", action)
	}

    // 3. Блок обработки ошибок (ДО запуска горутин!)
	if err != nil {
		h.Logger.Printf("Ошибка маршрутизации игрока %s: %v. Закрытие соединения.", playerID, err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"error": "Ошибка подключения: %s"}`, err.Error())))
		conn.Close()
		return // <-- Прерываем выполнение функции, горутины НЕ ЗАПУСТЯТСЯ
	}

	if err := h.sendRoomInfo(currentPlayer, currentRoom, messageType); err != nil {
		h.Logger.Printf("Ошибка отправки информации о комнате игроку %s: %v. Закрытие соединения.", currentPlayer.ID, err)
		// Если не удалось отправить первое сообщение, соединение лучше закрыть
		conn.Close()
		return
	}

	go currentPlayer.WritePump(cfg) 
	go currentPlayer.ReadPump(currentRoom, cfg) // Используем currentRoom
}

func(h *MainHub) sendRoomInfo(player *player.Player, currentRoom *room.Room, messageType string) error {
	playersList := make([]map[string]string, 0, len(currentRoom.Players))

	for _, player := range currentRoom.Players {
		playersList = append(playersList, map[string]string{
			"id": player.ID,
			"role": player.Role,
			"nickname": player.Nickname,

		})
	}

	response := map[string]interface{}{
		"type": messageType,
		"roomID": currentRoom.ID,
		"roomName": currentRoom.RoomName,
		"players": playersList,
		"role": player.Role,
		"currentPlayerID": player.ID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON, number %s", err)
	}

	select {
    case player.Send <- jsonResponse:
        // Успешно поместили сообщение в буфер
        return nil
    default:
        // Если канал заполнен, значит, клиент не успевает читать.
        // Это сигнал о том, что соединение не в порядке (обычно его нужно закрыть).
        return fmt.Errorf("канал отправки (player.Send) переполнен для игрока %s", player.ID)
    }
}

func (h *MainHub) ReconnectPlayer(newConn *websocket.Conn, roomID string, playerID string) (*player.Player, error)  {
	h.RoomsMutex.RLock() // Используем блокировку для чтения, чтобы безопасно читать из map
    room, found := h.Rooms[roomID]
    h.RoomsMutex.RUnlock()

	if !found {
		newConn.Close()
        return nil, fmt.Errorf("комната с ID %s не найдена", roomID)
	}

	room.PlayersMutex.RLock() 
    oldPlayer, found := room.Players[playerID]
    room.PlayersMutex.RUnlock()
	
	if !found {
		newConn.Close()
		return nil, fmt.Errorf("игрок c ID %s не найден в комнате %s", playerID, roomID)
	}

	oldPlayer.Mutex.Lock()
	defer oldPlayer.Mutex.Unlock()

	close(oldPlayer.Send)
	oldPlayer.Conn.Close() // Закрываем старое соединение

	select {
    case <-oldPlayer.Done:
        // Успех: Старые горутины завершились.
    case <-time.After(100 * time.Millisecond): // Таймаут на всякий случай
        h.Logger.Printf("Предупреждение: Read/Write Pump игрока %s не завершились за 3 секунды. Продолжаем с риском.", playerID)
    }

    oldPlayer.Conn = newConn // Присваиваем НОВОЕ соединение
	oldPlayer.Send = make(chan []byte, 256)
	oldPlayer.Done = make(chan struct{})

	room.BroadcastRoomUpdate() // <--- ДОБАВЛЕНО

	h.Logger.Printf("✅ Игрок %s успешно переподключен к комнате %s.", playerID, roomID)

	return oldPlayer, nil
}

func(h *MainHub) JoinRoom(player *player.Player, roomID string) (*room.Room, error) {
	h.RoomsMutex.RLock() // Используем блокировку для чтения, чтобы безопасно читать из map
    room, found := h.Rooms[roomID]
    h.RoomsMutex.RUnlock()

	if !found {return nil, fmt.Errorf("комната с ID %s не найдена", roomID)}
	
	room.PlayersMutex.Lock() 
	room.Players[player.ID] = player
	room.PlayersMutex.Unlock()

	room.BroadcastRoomUpdate()

	return room, nil
}