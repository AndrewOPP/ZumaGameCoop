export class NetworkManager {
    constructor(onStateReceived,) {
        this.ws = new WebSocket("ws://localhost:8080/ws");
        this.onStateReceived = onStateReceived; // Коллбэк для передачи данных в Renderer

        this.playerID = "CLIENT_A_123";
        this.ws.onopen = () => {
            console.log("WebSocket connection done");
            // this.ws.send("Hello GO! Connection test");
        };

        this.ws.onmessage = (event) => {
            try {
    
                const gameState = JSON.parse(event.data);
                // Передаем данные в внешний обработчик (Renderer)
                this.onStateReceived(gameState); 
            } catch (e) {
                console.error("Failed to parse game state:", e);
            }
        };

        this.ws.onerror = (error) => console.error("WebSocket error:", error);
        this.ws.onclose = () => console.log("Connection closed");
    }

    /**
     * Формирует и отправляет стандартизированную команду на сервер.
     * @param {string} commandType - Тип команды (например, 'move', 'shoot').
     * @param {Object} commandData - Данные, специфичные для команды (например, {x: 10, y: 5}).
     * @param {Object} payloadData - Данные, специфичные для команды (например, {x: 10, y: 5}).
     */
    sendCommand(commandType, commandData = {}, payloadData = {}) {
        if (this.ws.readyState !== WebSocket.OPEN) {
            console.warn(`WebSocket is not open. Command '${commandType}' ignored.`);
            return;
        }

        const commandPayload = {
            player_id: this.playerID, // Используем установленный ID
            type: commandType,
            data: commandData,
            payload: payloadData,
        };

        const jsonCommand = JSON.stringify(commandPayload);
        
        // Используем метод send() WebSocket API
        this.ws.send(jsonCommand);
        console.log(`Command sent: ${commandType}`, commandPayload);
    }

    send(message) {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(message);
        }
    }
}

