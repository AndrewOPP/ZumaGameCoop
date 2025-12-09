export class NetworkManager {
    // Явно объявляем свойства класса с типом any
    private ws: any;
    private onStateReceived: any;
    private playerID: string;

    /**
     * @param {any} onStateReceived - Коллбэк для передачи данных в Renderer (тип any)
     */
    constructor(onStateReceived: any) {
        // Проставляем типы в конструкторе
        this.ws = new WebSocket("ws://localhost:8080/ws");
        this.onStateReceived = onStateReceived; // Коллбэк для передачи данных в Renderer

        this.playerID = "CLIENT_A_123";
        this.ws.onopen = () => {
            console.log("WebSocket connection done");
            // this.ws.send("Hello GO! Connection test");
        };

        this.ws.onmessage = (event: any) => {
            try {
                // event.data неявно имеет тип any, если WebSocket не типизирован
                const gameState: any = JSON.parse(event.data);
                // Передаем данные в внешний обработчик (Renderer)
                this.onStateReceived(gameState); 
            } catch (e) {
                console.error("Failed to parse game state:", e);
            }
        };

        this.ws.onerror = (error: any) => console.error("WebSocket error:", error);
        this.ws.onclose = () => console.log("Connection closed");
    }

    /**
     * Формирует и отправляет стандартизированную команду на сервер.
     * @param {any} commandType - Тип команды (например, 'move', 'shoot').
     * @param {any} commandData - Данные, специфичные для команды (например, {x: 10, y: 5}).
     * @param {any} payloadData - Данные, специфичные для команды (например, {x: 10, y: 5}).
     */
    sendCommand(commandType: any, commandData: any = {}, payloadData: any = {}): void {
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

    /**
     * @param {any} message - Сообщение для отправки.
     */
    send(message: any): void {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(message);
        }
    }
}