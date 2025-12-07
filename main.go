package main

import (
	"encoding/json"
	"fmt"
	"github.com/AndrewOPP/ZumaGameCoop/config"
	"github.com/AndrewOPP/ZumaGameCoop/game"
	"github.com/AndrewOPP/ZumaGameCoop/hub"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	// Разрешаем все домены (для простого примера)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Апгрейд HTTP соединения в WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		log.Println("New WebSocket connection!")
		// conn.WriteMessage(websocket.TextMessage, []byte("Welcome! Connection established."))
		// Передаем conn и экземпляр Хаба в обработчик
		go handleMessages(conn, h)
	}
}

func handleMessages(conn *websocket.Conn, h *hub.Hub) {
	h.Register <- conn

	defer func() {
		h.Unregister <- conn
		log.Println("Connection closed and unregistered")
	}() // Гарантированное закрытие соединения

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received from client: %s\n", message)
		// Эхо обратно

		var cmd hub.PlayerCommand // Предполагаем, что PlayerCommand определен в hub
		err = json.Unmarshal(message, &cmd)

		if err == nil {
			// УСПЕШНАЯ ДЕСЕРИАЛИЗАЦИЯ: Это валидная команда.
			// log.Printf("Received command type: %s\n", cmd.CommandType)
			// log.Printf("Received input command: %+v", cmd)
			h.InputGate <- cmd
			continue // Переходим к следующей итерации цикла
		}

		// Отправляем структурированный объект команды (hub.PlayerCommand)
		// в InputGate, который слушает GameManager.
		// conn.WriteMessage(messageType, []byte("Server echoes: "+string(message)))
	}
}
func main() {
	cfg := config.LoadConfig()

	h := hub.NewHub()
	go h.Run()

	gm := game.NewGameManager(h, cfg)
	go gm.Run() // <-- Запуск в своей горутине!

	fmt.Println("Сервер запущен на " + cfg.Server.Host + cfg.Server.Port)

	fs := http.FileServer(http.Dir("./" + cfg.Server.FrontMainFolder))

	http.Handle("/", fs)
	http.HandleFunc("/ws", wsHandler(h))

	// 2. Запускаем сервер:
	// Слушаем порт 8080.
	// Если ListenAndServe вернет ошибку, паникуем и выводим ее.
	err := http.ListenAndServe(cfg.Server.Port, nil)
	if err != nil {
		// Обязательно проверяйте ошибку!
		panic(err)
	}
}
