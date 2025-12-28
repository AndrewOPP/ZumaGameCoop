// import GameRow from './GameRow/GameRow';
// import './Gameview.css';
// import { NetworkManager, type RoomInfo } from '../../game/network';
// import GameKeyboardView from './GameKeyboardView/GameKeyboardView';
// import { useEffect, useState } from 'react';

// interface GameViewProps {
//   networkManager: NetworkManager;
//   roomData: RoomInfo;
// }

// export default function GameView({ networkManager, roomData }: GameViewProps) {
//   const [inputValue, setInputValue] = useState('');

//   const {
//     gameState: { playerAttempts, scores },
//     currentPlayerID,
//   } = roomData;

//   const HandleOnSubmit = (word: string) => {
//     if (inputValue.length == 6) {
//       networkManager.sendCommand(roomData.currentPlayerID, 'check_word', { word: word.toLocaleUpperCase() });
//     }
//   };

//   // Внутри GameView
//   useEffect(() => {
//     const onKeyDown = (e: KeyboardEvent) => {
//       if (e.ctrlKey || e.metaKey || e.altKey) return;

//       // Вся логика теперь живет здесь
//       const key = e.key;
//       if (key === 'Backspace') {
//         setInputValue((prev) => prev.slice(0, -1));
//       } else if (key === 'Enter') {
//         setInputValue((prev) => {
//           if (prev.length === 5) {
//             HandleOnSubmit(prev); // Важно: HandleOnSubmit тоже должна быть стабильной
//             return '';
//           }
//           return prev;
//         });
//       } else if (/^[a-zA-Zа-яА-ЯёЁ]$/.test(key)) {
//         setInputValue((prev) => (prev.length < 5 ? prev + key.toUpperCase() : prev));
//       }
//     };

//     window.addEventListener('keydown', onKeyDown);
//     return () => window.removeEventListener('keydown', onKeyDown);
//   }, [HandleOnSubmit]); // Теперь только одна зависимость

//   const chooseValue = (rowIndex: number, currentRow: number) => {
//     if (rowIndex === currentRow) {
//       // console.log(currentRow, 'currentRow');
//       return inputValue;
//     }
//     return '';
//   };

//   const renderRows = () => {
//     const totalRows = 6;
//     const attempts = playerAttempts[currentPlayerID] || [];
//     return Array.from({ length: totalRows }).map((_, index) => {
//       return (
//         <GameRow
//           key={index}
//           onWordSubmit={HandleOnSubmit}
//           attempt={attempts[index]}
//           value={chooseValue(index, attempts.length)}
//           setInputValue={setInputValue}
//         />
//       );
//     });
//   };

//   return (
//     <div className="mainGameContainer">
//       <p>Конец игры через: {roomData.gameState.timeRemaining}</p>
//       <p>Score: {scores[currentPlayerID]}</p>
//       <div className="rowsContainer">{renderRows()}</div>
//       <GameKeyboardView roomData={roomData} networkManager={networkManager} setInputValue={setInputValue} />
//     </div>
//   );
// }

import { useEffect, useState, useCallback } from 'react';
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

  const {
    gameState: { playerAttempts, scores, timeRemaining },
    currentPlayerID,
  } = roomData;

  // 1. Оборачиваем отправку в useCallback, чтобы она не менялась при каждом рендере
  const HandleOnSubmit = useCallback(
    (word: string) => {
      // Важно: в Wordle длина слова обычно 5
      if (word.length === 5) {
        networkManager.sendCommand(currentPlayerID, 'check_word', {
          word: word.toUpperCase(),
        });
      }
    },
    [networkManager, currentPlayerID],
  );

  // 2. Единая логика обработки клавиш (и для физической, и для экранной клавы)
  const handleInput = useCallback(
    (key: string) => {
      if (key === 'Backspace' || key === '⌫') {
        setInputValue((prev) => prev.slice(0, -1));
      } else if (key === 'ENTER' || key === 'Enter') {
        if (!isValidWord(inputValue)) {
          console.log('Такого слова нет в словаре');
          // triggerShakeAnimation(); // Трясем инпут
          return; // Не отправляем на сервер мусор
        }
        HandleOnSubmit(inputValue);
        setInputValue('');
      } else if (/^[a-zA-Zа-яА-ЯёЁ]$/.test(key)) {
        setInputValue((prev) => (prev.length < 5 ? prev + key.toUpperCase() : prev));
      }
    },
    [HandleOnSubmit, inputValue],
  );

  // 3. Глобальный слушатель физической клавиатуры
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey || e.metaKey || e.altKey) return;
      handleInput(e.key);
    };

    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [handleInput]);

  const renderRows = () => {
    const totalRows = 6;
    const attempts = playerAttempts[currentPlayerID] || [];

    return Array.from({ length: totalRows }).map((_, index) => {
      // Определяем, какую строку передать в ряд
      const isCurrentRow = index === attempts.length;
      const rowValue = isCurrentRow ? inputValue : '';

      return <GameRow key={index} attempt={attempts[index]} value={rowValue} />;
    });
  };

  return (
    <div className="mainGameContainer">
      <p>Конец игры через: {timeRemaining}</p>
      <p>Score: {scores[currentPlayerID]}</p>
      <div className="rowsContainer">{renderRows()}</div>
      {/* 4. Передаем handleInput в экранную клавиатуру */}
      <GameKeyboardView roomData={roomData} networkManager={networkManager} onKeyPress={handleInput} />
    </div>
  );
}
