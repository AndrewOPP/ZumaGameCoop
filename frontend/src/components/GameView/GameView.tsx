import { useEffect, useState, useCallback, useMemo } from 'react';
import GameRow from './GameRow/GameRow';
import GameKeyboardView from './GameKeyboardView/GameKeyboardView';
import { NetworkManager, type RoomInfo } from '../../game/network';
import './Gameview.css';
import { isValidWord } from '../../utils/dictionary';

interface GameViewProps {
  networkManager: NetworkManager;
  roomData: RoomInfo;
}

export default function GameView({ networkManager, roomData }: GameViewProps) {
  const [inputValue, setInputValue] = useState('');
  const [isShaking, setIsShaking] = useState(false);

  // Добавляем rowIndex, чтобы слово не прыгало между строками одного раунда
  const [localAttempt, setLocalAttempt] = useState<{
    word: string;
    scoreAtTime: number;
    rowIndex: number;
  } | null>(null);

  // Безопасно извлекаем данные (защита от undefined в начале игры)
  const { gameState, currentPlayerID } = roomData;
  const player = gameState.players[currentPlayerID];
  const playerAttempts = gameState.playerAttempts[currentPlayerID] || [];

  const currentScore = player?.score || 0;
  const isServerWaiting = player?.isWaiting || false;
  const nextIndex = playerAttempts.length;

  // --- ВЫЧИСЛЯЕМОЕ СОСТОЯНИЕ (Derived State) ---
  // Слово валидно ТОЛЬКО если совпадает и счет (раунд), и номер строки
  const currentLocalWord = useMemo(() => {
    if (localAttempt && localAttempt.scoreAtTime === currentScore && localAttempt.rowIndex === nextIndex) {
      return localAttempt.word;
    }
    return null;
  }, [localAttempt, currentScore, nextIndex]);

  const triggerShake = () => {
    setIsShaking(true);
    setTimeout(() => setIsShaking(false), 700);
  };

  const HandleOnSubmit = useCallback(
    (word: string) => {
      if (word.length === 5) {
        const formatted = word.toUpperCase();

        // Фиксируем слово за конкретной строкой конкретного раунда
        setLocalAttempt({
          word: formatted,
          scoreAtTime: currentScore,
          rowIndex: nextIndex,
        });

        networkManager.sendCommand(currentPlayerID, 'check_word', { word: formatted });
        setInputValue('');
      }
    },
    [networkManager, currentPlayerID, currentScore, nextIndex],
  );

  const handleInput = useCallback(
    (key: string) => {
      // Если сервер "думает", блокируем только ENTER и ввод новых букв,
      // но даем стирать (Backspace), если это нужно по логике игры.
      if (isServerWaiting) return;

      if (key === 'Backspace' || key === '⌫') {
        setInputValue((prev) => prev.slice(0, -1));
      } else if (key === 'ENTER' || key === 'Enter') {
        if (inputValue.length < 5 || !isValidWord(inputValue)) {
          triggerShake();
          return;
        }
        HandleOnSubmit(inputValue);
      } else if (/^[a-zA-Zа-яА-ЯёЁ]$/.test(key)) {
        setInputValue((prev) => (prev.length < 5 ? prev + key.toUpperCase() : prev));
      }
    },
    [HandleOnSubmit, inputValue, isServerWaiting],
  );

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey || e.metaKey || e.altKey) return;
      handleInput(e.key);
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [handleInput]);

  const renderRows = () => {
    return Array.from({ length: 6 }).map((_, index) => {
      const serverAttempt = playerAttempts[index];
      const isCurrentRow = index === nextIndex;

      if (serverAttempt) {
        return <GameRow key={`done-${index}`} attempt={serverAttempt} value={serverAttempt.word} />;
      }

      // Если это текущая строка: приоритет у "зависшего" слова, потом ввод.
      // Если строка не текущая — всегда пусто.
      const rowValue = isCurrentRow ? currentLocalWord || inputValue : '';

      return (
        <div key={`active-${index}`} className={isCurrentRow && isShaking ? 'shake' : ''}>
          <GameRow value={rowValue} />
        </div>
      );
    });
  };

  return (
    <div className="mainGameContainer">
      <div className="gameHeader">
        <p>Конец игры через: {gameState.timeRemaining}</p>
        <p>Score: {currentScore}</p>
      </div>
      <div className="rowsContainer">{renderRows()}</div>
      <GameKeyboardView roomData={roomData} networkManager={networkManager} onKeyPress={handleInput} />
    </div>
  );
}
