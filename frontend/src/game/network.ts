const SERVER_ADDRESS = 'ws://localhost:8080/ws';

export interface WordleAttempt {
  word: string;
  result: string;
  isCorrect?: boolean;
}

export interface Player {
  id: string;
  nickname: string;
  score: number;
  role: string;
  isWaiting: boolean; // Тот самый флаг для 2-секундной паузы
  isReady: boolean;
  roomId: string;
}

export interface GameState {
  // map[string]string -> Record<string, string>
  currentWords: Record<string, string>;

  // map[string][]WordleAttempt -> Record<string, WordleAttempt[]>
  playerAttempts: Record<string, WordleAttempt[]>;

  // map[string]int -> Record<string, number>
  scores: Record<string, number>;

  // map[string]*Player -> Record<string, Player>
  // Теперь это мапа (ID -> Объект игрока), как мы и договорились на бэкенде
  players: Record<string, Player>;

  timeRemaining: number;
  isActive: boolean;
  readyStatus: Record<string, boolean>;
}

// Интерфейс для структурирования данных, которые мы ожидаем от сервера
export interface RoomInfo {
  type: 'room_created' | 'room_joined' | 'room_reconnected' | 'error' | 'STATUS' | 'room_updated';
  roomID: string;
  role: string;
  roomName: string;
  error?: string;
  currentPlayerID: string;
  gameState: GameState;
}

// Интерфейс команды, которую мы отправляем на сервер (соответствует Go PlayerRawCommand)
export interface Command {
  playerID: string;
  type: string;
  data?: any; // Данные, специфичные для команды (например, { angle: 45 })
}

export class NetworkManager {
  private ws: WebSocket | null = null;

  // Коллбэк теперь типизирован, чтобы принимать структурированные данные о комнате/игре
  private onMessageReceived: (data: any) => void;

  // Новые поля для хранения информации после успешного подключения
  public playerID: string = '';
  public roomId: string = '';

  constructor(onMessageReceied: (data: any) => void) {
    this.onMessageReceived = onMessageReceied;
  }

  /**
   * Устанавливает соединение с сервером для создания или присоединения к комнате.
   * @param {string} action - 'create' или 'join'.
   * @param {string} playerID - Имя игрока.
   * @param {string | null} [roomId=null] - ID комнаты, если action === 'join'.
   */
  public connect(
    action: 'create' | 'join' | 'reconnect',
    playerID: string | null,
    roomId: string | null = null,
    nickname: string,
    roomName?: string,
  ): void {
    if (this.ws) {
      console.warn('Connection already established.');
      return;
    }

    const url = new URL(SERVER_ADDRESS);
    url.searchParams.append('action', action);
    url.searchParams.append('nickname', nickname);

    if (playerID) url.searchParams.append('playerID', playerID);
    if (roomId) url.searchParams.append('roomId', roomId);
    if (roomName) url.searchParams.append('roomName', roomName);

    const ws = new WebSocket(url.toString());
    this.ws = ws;

    ws.onopen = () => {
      console.log('✅ WebSocket connected for action: ${action}');
      this.onMessageReceived({
        type: 'STATUS',
        message: 'Connected to server.',
      });
    };

    ws.onmessage = (event) => {
      try {
        const data: RoomInfo = JSON.parse(event.data);
        // console.log(data, 'datadatadatadatadata');

        if (['room_created', 'room_joined', 'room_updated', 'room_reconnected'].includes(data.type)) {
          this.roomId = data.roomID;

          // Используем currentPlayerID, который сервер прислал в корне объекта
          if (data.currentPlayerID) {
            this.playerID = data.currentPlayerID;
          }
        }
        this.onMessageReceived(data); // Передаем данные в React-компонент
      } catch (error) {
        console.error('Failed to parse message from server:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('⚠️ WebSocket error:', error);
      this.onMessageReceived({ type: 'STATUS', message: 'Connection Error.' });
    };

    ws.onclose = () => {
      console.log('Connection closed.');
      this.onMessageReceived({ type: 'STATUS', message: 'Disconnected.' });
      this.ws = null;
    };
  }

  /**
   * Формирует и отправляет стандартизированную команду на сервер.
   * @param {string} type - Тип команды (например, 'move_cannon', 'fire').
   * @param {any} data - Данные, специфичные для команды (соответствует Go rawCmd.Data).
   */
  public sendCommand(playerID: string, type: string, data: any = {}): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn(`WebSocket is not open. Command '${type}' ignored.`);
      return;
    }

    // 3. Упрощенный Payload, соответствующий Go PlayerRawCommand
    const commandPayload: Command = {
      type: type,
      data: data,
      playerID: playerID, // Объект или Map, который будет сериализован в JSON
    };

    const jsonCommand = JSON.stringify(commandPayload);
    this.ws.send(jsonCommand);
    console.log(`Command sent: ${type}`, commandPayload);
  }
}
