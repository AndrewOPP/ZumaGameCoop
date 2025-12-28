import type { NetworkManager, RoomInfo } from '../../../game/network';
import './GameKeyboardView.css';

const ROWS = [
  ['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'],
  ['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'],
  ['ENTER', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', '⌫'],
];

interface GameKeyboardProps {
  roomData: RoomInfo; // Ожидаем объект roomData
  networkManager: NetworkManager; // И объект networkManager
  onKeyPress: (key: string) => void;
}

export default function GameKeyboardView({ roomData, networkManager, onKeyPress }: GameKeyboardProps) {
  const currentPlayersAttempts = roomData.gameState.playerAttempts[roomData.currentPlayerID];

  const getKeyStatusClass = (status: string): string => {
    if (status === 'G') return 'key-correct';
    if (status === 'Y') return 'key-present';
    if (status === 'X') return 'key-absent';
    return ''; // Если буквы еще не было в попытках
  };

  const STATUS_WEIGHTS: Record<string, number> = {
    G: 3, // Самый приоритетный (Зеленый)
    Y: 2, // Средний (Желтый)
    X: 1, // Низкий (Серый)
    '': 0, // Еще не нажимали
  };

  const allPlayerCharsAndResults = () => {
    const resultsObj: Record<string, string> = {};
    let playerUsedChars = ''; // Используем строку вместо массива для простоты

    for (const attempt of currentPlayersAttempts) {
      playerUsedChars += attempt.word;

      for (let j = 0; j < attempt.word.length; j++) {
        const char = attempt.word[j];
        const newStatus = attempt.result[j];

        // Достаем вес текущего статуса буквы (если буквы нет, вес будет 0)
        const currentWeight = STATUS_WEIGHTS[resultsObj[char]] || 0;
        const newWeight = STATUS_WEIGHTS[newStatus] || 0;

        // ЗАПИСЫВАЕМ ТОЛЬКО ЕСЛИ НОВЫЙ СТАТУС "ТЯЖЕЛЕЕ"
        if (newWeight > currentWeight) {
          resultsObj[char] = newStatus;
        }
      }
    }

    return { playerUsedChars, resultsObj };
  };

  const playerUsedInfo = allPlayerCharsAndResults();

  return (
    <div className="keyboard-container">
      {ROWS.map((row, rowIndex) => (
        <div key={rowIndex} className="keyboard-row">
          {row.map((key) => {
            // Определяем класс для специальных кнопок
            // const isSpecial = key === 'EER' || key === '⌫';
            if (playerUsedInfo?.playerUsedChars.includes(key)) {
              return (
                <button
                  onClick={(e) => onKeyPress(e.currentTarget.textContent || '')}
                  key={key}
                  className={`keyboard-key ${getKeyStatusClass(playerUsedInfo.resultsObj[key].toUpperCase())}`}
                >
                  {key}
                </button>
              );
            }

            return (
              <button
                onClick={(e) => onKeyPress(e.currentTarget.textContent || '')}
                onMouseDown={(e) => e.preventDefault()} // Чтобы не терять фокус, если он где-то остался
                key={key}
                className={`keyboard-key`}
              >
                {key}
              </button>
            );
          })}
        </div>
      ))}
    </div>
  );
}
