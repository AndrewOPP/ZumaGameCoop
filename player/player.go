package player

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/AndrewOPP/ZumaGameCoop/config"
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
		Done: make(chan struct{}),
	}
}

func (player *Player) WritePump(cfg *config.Config) {

	ticker := time.NewTicker(cfg.WebSocket.PingPeriod)

	defer func() {
		// player.Conn.Close()
	}()

	for {

		select {
		// 3. Обработка исходящих игровых сообщений
		case message, ok := <-player.Send:
			// Устанавливаем таймаут для этой операции записи (чтобы не зависло)
			player.Conn.SetWriteDeadline(time.Now().Add(cfg.WebSocket.WriteWait))
            
			if !ok {
				// Хаб закрыл канал, отправляем клиенту сообщение о закрытии и выходим
				player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Отправка игрового сообщения
			if err := player.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Write error (игровое сообщение) для %s: %v", player.ID, err)
				return // Возвращаемся при ошибке записи
			}

		// 4. Обработка тикера (Отправка Ping)
		case <-ticker.C:
			// Устанавливаем таймаут для этой операции записи Ping
			player.Conn.SetWriteDeadline(time.Now().Add(cfg.WebSocket.WriteWait)) 
            
			// Отправляем Ping-сообщение (data = nil)
			if err := player.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Write error (Ping) для %s: %v", player.ID, err)
				return // Возвращаемся при ошибке, если не удалось отправить Ping
			}
		}
}
}

func (player *Player) ReadPump(room RoomContext, cfg *config.Config) {

	player.Conn.SetReadLimit(cfg.WebSocket.MaxMessageSize)
	player.Conn.SetReadDeadline(time.Now().Add(cfg.WebSocket.PongWait))

	player.Conn.SetPongHandler(func(string) error {
		player.Conn.SetReadDeadline((time.Now().Add(cfg.WebSocket.PongWait)))
		return nil
	})


	defer func() {
		// Уведомляем комнату/хаб об отключении (если есть канал Unregister)
        // Если в комнате есть канал Unregister chan *Player:
        // r.Unregister <- p 
        
        // Закрываем соединение WebSocket
		close(player.Done)
		// player.Conn.Close()
		log.Printf("ReadPump завершен для игрока %s", player.ID)
	}()

		for {
			_, message, err := player.Conn.ReadMessage()

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("ReadPump: Неожиданное закрытие соединения от %s: %v", player.ID, err)
				} else {
					// Это может быть таймаут, закрытие, или другая ожидаемая ошибка
					log.Printf("ReadPump: Игрок %s отключается (ошибка чтения): %v", player.ID, err)
				}
				return
			}

			var rawCmd PlayerRawCommand
			if err := json.Unmarshal(message, &rawCmd); err != nil {
				log.Printf("Ошибка десериализации команды от %s: %v. Сырое сообщение: %s", player.ID, err, message)
            continue // Игнорируем невалидный JSON и ждем следующего сообщения
			}

			commandForRoom := &PlayerCommand{
				PlayerID: player.ID,
				Type: rawCmd.Type,
				Data: rawCmd.Data,
			}

			select {
			case room.InputGateChan() <- commandForRoom:
				log.Printf("Команда %s от игрока %s успешно передана в комнату.", commandForRoom.Type, player.ID)
			default:
				log.Printf("InputGate комнаты %s переполнен, команда %s пропущена.", room.GetID(), commandForRoom.Type)
			}

		}
}
