package main

import (
	// "encoding/json"
	"fmt"
	"github.com/AndrewOPP/ZumaGameCoop/config"
	// "github.com/AndrewOPP/ZumaGameCoop/game"
	"github.com/AndrewOPP/ZumaGameCoop/mainHub"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var upgrader = websocket.Upgrader{
	// Разрешаем все домены (для простого примера)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(h *mainhub.MainHub) http.HandlerFunc {
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
		// go handleMessages(conn, h)
		go h.RoutePlayer(conn, r)
	}
}

// func handleMessages(conn *websocket.Conn, h *mainhub.MainHub) {
// 	// h.Register <- conn

// 	defer func() {
// 		// h.Unregister <- conn
// 		log.Println("Connection closed and unregistered")
// 	}() // Гарантированное закрытие соединения

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			break
// 		}
// 		log.Printf("Received from client: %s\n", message)
// 		// Эхо обратно

// 		// var cmd hub.PlayerCommand // Предполагаем, что PlayerCommand определен в hub
// 		// err = json.Unmarshal(message, &cmd)

// 		if err == nil {
// 			// УСПЕШНАЯ ДЕСЕРИАЛИЗАЦИЯ: Это валидная команда.
// 			// log.Printf("Received command type: %s\n", cmd.CommandType)
// 			// log.Printf("Received input command: %+v", cmd)
// 			// h.InputGate <- cmd
// 			continue // Переходим к следующей итерации цикла
// 		}

// 		// Отправляем структурированный объект команды (hub.PlayerCommand)
// 		// в InputGate, который слушает GameManager.
// 		// conn.WriteMessage(messageType, []byte("Server echoes: "+string(message)))
// 	}
// }

func spaHandler(buildPath string) http.HandlerFunc {
	// Создаем FileServer для обслуживания статических файлов
	fs := http.FileServer(http.Dir(buildPath))

	return func(w http.ResponseWriter, r *http.Request) {
		// Формируем полный путь к запрошенному файлу
		filePath := buildPath + r.URL.Path

		// Проверяем, существует ли файл в файловой системе
		_, err := os.Stat(filePath)

		// Если файл НЕ существует (os.IsNotExist) или произошла другая ошибка,
		// это, вероятно, роут клиента (SPA). Возвращаем index.html.
		if os.IsNotExist(err) || err != nil {
			// log.Printf("File not found at %s. Serving index.html (SPA Fallback).", filePath)
			http.ServeFile(w, r, buildPath+"/index.html")
			return
		}

		// Если файл существует, отдаем его с помощью FileServer.
		fs.ServeHTTP(w, r)
	}
}

func main() {
	cfg := config.LoadConfig()

	isDevMode := os.Getenv("DEV_MODE") == "true"
	const reactDevServerURL = "http://localhost:5173"
	const reactBuildPath = "frontend/dist"

	h := mainhub.NewMainHub()
	// go h.Run()

	if isDevMode {
		fmt.Println("🚀 Включен режим разработки. Фронтенд проксируется на", reactDevServerURL)

		// Создаем целевой URL для прокси
		proxyURL, _ := url.Parse(reactDevServerURL)

		// Создаем Reverse Proxy
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)

		// Регистрируем обработчик, который проксирует запросы на Dev Server
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Важное исключение: если это не API и не WS, проксируем.
			// Если у вас есть другие API, добавьте исключения здесь.
			if r.URL.Path == "/ws" {
				// Это должно быть обработано выше, но как защита:
				wsHandler(h).ServeHTTP(w, r)
				return
			}
			// Проксируем все остальные запросы на React Dev Server (localhost:5173)
			proxy.ServeHTTP(w, r)
		})

	} else {
		// Режим продакшена: используем статические файлы и SPA Fallback
		fmt.Println("📦 Включен режим продакшена. Обслуживание статических файлов.")
		http.HandleFunc("/", spaHandler(reactBuildPath))
	}

	// gm := game.NewGameManager(h, cfg)
	// go gm.Run()

	fmt.Println("Сервер запущен на " + cfg.Server.Host + cfg.Server.Port)

	// 1. Регистрируем обработчик для WebSocket
	// http.HandleFunc("/ws", wsHandler(h))

	// // 2. Регистрируем обработчик для фронтенда (все остальные запросы)
	// // Передаем путь к билду в функцию, которая вернет обработчик.
	// http.HandleFunc("/", spaHandler(reactBuildPath))

	// 3. Запускаем сервер
	err := http.ListenAndServe(cfg.Server.Port, nil)
	if err != nil {
		panic(err)
	}
}
