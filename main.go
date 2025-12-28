package main

import (
	// "encoding/json"
	"fmt"
	"github.com/AndrewOPP/ZumaGameCoop/config"
	"github.com/AndrewOPP/ZumaGameCoop/mainHub"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
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
    const reactBuildPath = "frontend/dist"
    
    h := mainhub.NewMainHub(cfg)

    // 1. WebSocket регистрируем ВСЕГДА и ПЕРВЫМ. 
    // Это отдельный роут, который не должен пересекаться со статикой.
    http.HandleFunc("/ws", wsHandler(h))

    if isDevMode {
        const reactDevServerURL = "http://localhost:5173"
        fmt.Println("🚀 Режим разработки: прокси на", reactDevServerURL)
        
        proxyURL, _ := url.Parse(reactDevServerURL)
        proxy := httputil.NewSingleHostReverseProxy(proxyURL)

        // В деве всё, кроме /ws, проксируем
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            proxy.ServeHTTP(w, r)
        })
    } else {
        fmt.Println("📦 Режим продакшена: раздача из", reactBuildPath)
        
        // Используем стандартный обработчик статики для существующих файлов
        fs := http.FileServer(http.Dir(reactBuildPath))

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            // Проверяем, существует ли физический файл (JS, CSS, картинка)
            // Используем filepath.Join для корректных путей в Windows/Linux
            path := filepath.Join(reactBuildPath, r.URL.Path)
            info, err := os.Stat(path)
            
            // Если файла нет или это папка — отдаем index.html (SPA Fallback)
            if os.IsNotExist(err) || info.IsDir() {
                http.ServeFile(w, r, filepath.Join(reactBuildPath, "index.html"))
                return
            }
            
            fs.ServeHTTP(w, r)
        })
    }

    fmt.Printf("🌍 Сервер: http://localhost%s\n", cfg.Server.Port)
    
    // Важно: проверь, чтобы cfg.Server.Port начинался с двоеточия, например ":8080"
    if err := http.ListenAndServe(cfg.Server.Port, nil); err != nil {
        log.Fatal("ListenAndServe Error:", err)
    }
}