import React, { useRef } from 'react';
import './GameRow.css';

interface GameRowProps {
  onWordSubmit: (word: string) => void;
  attempt: {
    word: string;
    result: string;
  };
}

export default function GameRow({ onWordSubmit, attempt }: GameRowProps) {
  const inputsRef = useRef<(HTMLInputElement | null)[]>([]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>, index: number) => {
    const { value } = e.target;

    if (value && index < 4) {
      inputsRef.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>, index: number) => {
    if (index > 0 && e.key === 'Backspace' && !e.currentTarget.value) {
      inputsRef.current[index - 1]?.focus();
    }

    if (e.key === 'Enter') {
      const resultWord = inputsRef.current.map((input) => input?.value || '').join('');

      if (resultWord.length === 5) {
        // console.log(resultWord);

        onWordSubmit(resultWord);
      }
    }
  };

  return (
    <form className="eachRowContainer">
      {Array.from({ length: 5 }).map((_, index) => {
        let backgroundColor;
        let charColor;
        const charResultColor: string = attempt?.result[index];
        charColor = 'black';
        backgroundColor = 'transparent';

        if (charResultColor === 'G') backgroundColor = '#6aaa64';
        if (charResultColor === 'Y') backgroundColor = '#c9b458';
        if (charResultColor === 'X') backgroundColor = '#787c7e';
        if (charResultColor) charColor = '#ffffff';

        return (
          <input
            style={{ backgroundColor: backgroundColor, color: charColor }}
            type="text"
            className="gameInput"
            maxLength={1}
            disabled={charResultColor ? true : false}
            onChange={(e) => handleChange(e, index)}
            onKeyDown={(e) => handleKeyDown(e, index)}
            key={index}
            ref={(el) => {
              inputsRef.current[index] = el;
            }}
          />
        );
      })}
    </form>
  );
}
