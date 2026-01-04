import { useEffect, useRef, useState } from 'react';
import './GameRow.css';

interface GameRowProps {
  attempt?: {
    word: string;
    result: string;
    isCorrect?: boolean;
  };
  value: string;
  isWinner?: boolean;
  isNeedToClear?: boolean; // Новый пропс для анимации удаления
}

export default function GameRow({ attempt, value, isNeedToClear }: GameRowProps) {
  const lastPlayedWord = useRef<string | null>(null);
  const [showClear, setShowClear] = useState(false);

  useEffect(() => {
    if (attempt && attempt.word !== lastPlayedWord.current) {
      lastPlayedWord.current = attempt.word;

      for (let i = 0; i < 5; i++) {
        setTimeout(() => {
          const audio = new Audio('/audio/char_reveal.mp3');
          audio.volume = 0.15;
          audio.play().catch(() => {});
        }, i * 230);
      }
    }
  }, [attempt]);

  // Задержка перед очисткой после победы
  useEffect(() => {
    if (isNeedToClear) {
      const timer = setTimeout(() => {
        setShowClear(true);
      }, 3000); // 3 секунды
      return () => clearTimeout(timer);
    }
  }, [isNeedToClear]);

  const clearClass = showClear ? 'clear-after-win' : '';
  const isWinner = attempt?.isCorrect ? 'winner' : '';

  return (
    <div className="eachRowContainer">
      {Array.from({ length: 5 }).map((_, index) => {
        const charResultColor = attempt?.result[index];
        const charToShow = attempt ? attempt.word[index] : value[index] || '';
        const isFilled = !attempt && value[index];
        let colorClass = '';
        if (charResultColor === 'G') colorClass = 'correct';
        else if (charResultColor === 'Y') colorClass = 'present';
        else if (charResultColor === 'X') colorClass = 'absent';

        return (
          <div
            key={index}
            className={`gameCell cell-${index}
            ${isFilled ? 'filled' : ''}
            ${attempt ? 'flip' : ''}
            ${colorClass}
            ${isWinner}
            ${clearClass}`}
          >
            <span className="cell-text">{charToShow}</span>
          </div>
        );
      })}
    </div>
  );
}
