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
	"github.com/AndrewOPP/ZumaGameCoop/wordsmap"
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
	dictMap, dictSlice := wordsmap.LoadEmbeddedDictionary()

	// 2. Создаем уникальный ID комнаты (например, с помощью библиотеки uuid или просто счетчика)
    // Для простоты используем текущее время и имя хоста для псевдо-уникальности
	roomID := fmt.Sprintf("room_%d", time.Now().UnixNano())



	newRoom := &room.Room{
		ID: roomID,
		RoomName: roomName,
		HostID: HostConnection.ID,
		Players: make(map[string]*player.Player),
		StateUpdates: make(chan []byte),
		InputGate: make(chan *player.PlayerCommand),
		State: NewGameState(),
		Dictionary: dictMap,   // Мапа для проверок
		WordList:   dictSlice,
	}

	newRoom.Players[HostConnection.ID] = HostConnection

	h.Rooms[roomID] = newRoom

	h.Logger.Printf("Новая комната создана. ID: %s, Хост: %s", roomID, HostConnection.ID)

	return newRoom, nil
}

func NewGameState()  room.GameState {
    return  room.GameState{
        CurrentWords:  make(map[string]string),
		PlayerAttempts: make(map[string][]room.WordleAttempt),
		Players: make(map[string]*player.Player),
        TimeRemaining: 600,
        IsActive:      false,

    }
}



func (h *MainHub) RoutePlayer(conn *websocket.Conn, r *http.Request) {
    // Получаем доступ к конфигу через поле структуры MainHub (предполагаем, что h.Config существует)
    cfg := h.Config
	// player := player.CreatePlayer(conn, r)
	playerID := r.URL.Query().Get("playerID")
	action := r.URL.Query().Get("action")
	roomID := r.URL.Query().Get("roomId")
	roomName := r.URL.Query().Get("roomName")

	var err error
	// var messageType string
	// var currentPlayer *player.Player
	var result *RouteResult

	h.Logger.Printf("Входящее подключение. Действие: %s, Комната ID: %s, Игрок ID: %s", action, roomID, playerID)

	switch action {
	case "create":
		result, err = h.handleCreate(conn, r, roomName)
	case "join":
		result, err = h.handleJoin(conn, r)
	case "reconnect":
		result, err = h.handleReconect(conn, playerID, roomID)
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

	if err := h.sendRoomInfo(result.CurrentPlayer, result.CurrentRoom, result.MessageType); err != nil {
		h.Logger.Printf("Ошибка отправки информации о комнате игроку %s: %v. Закрытие соединения.", result.CurrentPlayer.ID, err)
		// Если не удалось отправить первое сообщение, соединение лучше закрыть
		conn.Close()
		return
	}	
	if action == "create" {
        go result.CurrentRoom.Run()
    }

	go result.CurrentPlayer.WritePump(cfg) 
	go result.CurrentPlayer.ReadPump(result.CurrentRoom, cfg) // Используем currentRoom
	// result.CurrentRoom.Run()
}

func(h *MainHub) sendRoomInfo(player *player.Player, currentRoom *room.Room, messageType string) error {
	currentRoom.PlayersMutex.RLock()
    defer currentRoom.PlayersMutex.RUnlock()
    
    currentRoom.Mu.Lock()
    defer currentRoom.Mu.Unlock()

	// playersList := make([]map[string]interface{}, 0, len(currentRoom.Players))

	// for _, player := range currentRoom.Players {
	// 	playersList = append(playersList, map[string]interface{}{
	// 		"id": player.ID,
	// 		"role": player.Role,
	// 		"nickname": player.Nickname,
	// 		"score": currentRoom.State.Players[player.ID].Score,
	// 		"isReady": currentRoom.State.Players[player.ID].IsReady,
	// 		"roomID": currentRoom.State.Players[player.ID].RoomId,

	// 	})
	// }

	response := map[string]interface{}{
		"type": messageType,
		"roomID": currentRoom.ID,
		"roomName": currentRoom.RoomName,
		// "players": playersList,
		// "role": player.Role, // убрать в будущем так как в плеер листе у нас есть роль для каждогл
		"currentPlayerID": player.ID,
		"gameState": currentRoom.State, // ПЕРЕДЕЛАТЬ, НЕ ПЕРЕДОВАТЬ ГОТОВНОСТЬ И СЛОВА ЧТОБЫ НЕ ЧИТЕРИЛИ 
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
        return nil, fmt.Errorf("комната c ID %s не найдена", roomID)
	}

	room.PlayersMutex.RLock() 
    // oldPlayer, found := room.Players[playerID]
    oldPlayer, found := room.State.Players[playerID]
    room.PlayersMutex.RUnlock()

	// НУЖНО НАХОДИТЬ В СТАРОЙ РУМЕ ПО АЙДИ ИГРОКА ЕГО ДАННЫЕ (ОЧКИ, ВРЕМЯ И ТД И ТУТ ПРИСВАИВАТЬ ЧТОБЫ НЕ БЫЛО ОШИБКИ)
	
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
	fmt.Printf("found %s", found)
	if !found {return nil, fmt.Errorf("комната с ID %s не найдена", roomID)}
	
	room.PlayersMutex.Lock() 
	room.State.Players[player.ID] = player
	room.Players[player.ID] = player
	room.PlayersMutex.Unlock()

	room.BroadcastRoomUpdate()

	return room, nil
}