package mainhub

import (
	"fmt"
	"net/http"
	"github.com/AndrewOPP/ZumaGameCoop/player"
	"github.com/AndrewOPP/ZumaGameCoop/room"
	"github.com/gorilla/websocket"
)

type RouteResult struct {
    CurrentPlayer *player.Player
    CurrentRoom   *room.Room // Предполагая, что room.Room импортирован или доступен
    MessageType   string
}

func (h *MainHub) handleJoin(conn *websocket.Conn, r *http.Request) (*RouteResult, error) {
	currentPlayer := player.CreatePlayer(conn, r) // conn должна быть доступна
	roomID := r.URL.Query().Get("roomId")

	currentPlayer.Role = "guest"
	currentRoom, err := h.JoinRoom(currentPlayer, roomID)
	fmt.Printf("currentRoom %s \n", currentRoom)
	fmt.Printf("err %s \n", err)
		if err != nil {
		return nil, fmt.Errorf("неудалось присоединиться: %w", err) // %w для обертывания ошибки
	}


	currentRoom.State.ReadyStatus[currentPlayer.ID] = false
	currentRoom.State.Scores[currentPlayer.ID] = 0
	currentRoom.State.PlayerAttempts[currentPlayer.ID] = make([]room.WordleAttempt, 0)


	h.Logger.Printf("Игрок %s присоединился к комнате %s.", currentPlayer.ID, roomID)
	return &RouteResult{
		CurrentPlayer: currentPlayer,
		CurrentRoom:   currentRoom,
		MessageType:   "room_joined",
	}, nil
}

func (h *MainHub) handleCreate(conn *websocket.Conn, r *http.Request, roomName string) (*RouteResult, error) {
	currentPlayer := player.CreatePlayer(conn, r)
	
	h.Logger.Printf("Игрок %s просит создать новую комнату.", currentPlayer.ID)

	currentPlayer.Role = "host"
	currentRoom, err := h.CreateRoom(currentPlayer, roomName)
	currentRoom.State.ReadyStatus[currentPlayer.ID] = false
	currentRoom.State.Scores[currentPlayer.ID] = 0
	currentRoom.State.PlayerAttempts[currentPlayer.ID] = make([]room.WordleAttempt, 0)
	
	if(err != nil) {
		h.Logger.Printf("Ошибка при создании комнаты. %s", err)
		return nil, fmt.Errorf("не удалось создать комнату: %w", err)
	}

	return &RouteResult{
		CurrentPlayer: currentPlayer,
		CurrentRoom:   currentRoom,
		MessageType:   "room_created",
	}, nil
}

func(h *MainHub) handleReconect(conn *websocket.Conn, playerID string, roomID string) (*RouteResult, error) {
	if(playerID == "" || roomID == "") {
		h.Logger.Println("Игрок неверен или комната не найдена")
		return nil, fmt.Errorf("игрок неверен или комната не найдена")
	}
	
	h.Logger.Printf("Игрок %s просит переподключение к комнате %s.", playerID, roomID)
	currentPlayer, err := h.ReconnectPlayer(conn, roomID, playerID)

	var currentRoom *room.Room
	
	if err == nil && currentPlayer != nil {
			h.mu.RLock()
			currentRoom = h.Rooms[roomID] 
			h.mu.RUnlock()
		}

	return &RouteResult{
		CurrentPlayer: currentPlayer,
		CurrentRoom:   currentRoom,
		MessageType:   "room_reconnected",
	}, nil
}

