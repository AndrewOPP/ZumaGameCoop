import { useEffect, useRef, useState } from 'react';

import './GameRow.css';

interface GameRowProps {
  // onWordSubmit здесь больше не нужен для работы инпутов,

  // но если он используется где-то еще, можно оставить.

  attempt?: {
    word: string;

    result: string;
  };

  value: string;
  isWinner?: boolean;
}

export default function GameRow({ attempt, value, isWinner }: GameRowProps) {
  // Храним последнее "озвученное" слово, чтобы звук не дублировался
  const lastPlayedWord = useRef<string | null>(null);

  useEffect(() => {
    // Проверяем:
    // 1. Есть ли попытка
    // 2. Не совпадает ли она с тем, что мы уже озвучили
    if (attempt && attempt.word !== lastPlayedWord.current) {
      // Сразу обновляем значение в Ref, чтобы при следующем тике условие не сработало
      lastPlayedWord.current = attempt.word;

      // Запускаем звуковую волну
      for (let i = 0; i < 5; i++) {
        setTimeout(() => {
          const audio = new Audio('/audio/char_reveal.mp3');
          audio.volume = 0.15;
          audio.play().catch(() => {
            // Браузер может блокировать звук до первого клика, это нормально
          });
        }, i * 230);
      }
    }
  }, [attempt]); // Эффект следит только за изменением пропса attempt

  return (
    <div className="eachRowContainer">
      {Array.from({ length: 5 }).map((_, index) => {
        const charResultColor = attempt?.result[index];
        const charToShow = attempt ? attempt.word[index] : value[index] || '';

        let colorClass = '';
        if (charResultColor === 'G') colorClass = 'correct';
        else if (charResultColor === 'Y') colorClass = 'present';
        else if (charResultColor === 'X') colorClass = 'absent';

        return (
          <div
            key={index}
            className={`gameCell cell-${index} ${attempt ? 'flip' : ''} ${colorClass} ${isWinner ? 'winner' : ''}`}
            // Важно: убираем borderColor из style, если есть attempt,
            // иначе он перекроет анимацию reveal-correct
            style={{
              borderColor: !attempt && value[index] ? '#878a8c' : undefined,
            }}
          >
            {charToShow}
          </div>
        );
      })}
    </div>
  );
}
