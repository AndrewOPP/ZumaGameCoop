// import React, { useEffect, useRef, type Dispatch, type SetStateAction } from 'react';
// import './GameRow.css';

// interface GameRowProps {
//   onWordSubmit: () => void;
//   attempt: {
//     word: string;
//     result: string;
//   };
//   value: string;
//   setInputValue: Dispatch<SetStateAction<string>>;
// }

// export default function GameRow({ onWordSubmit, attempt, value, setInputValue }: GameRowProps) {
//   const inputsRef = useRef<(HTMLInputElement | null)[]>([]);

//   const handleChange = (e: React.ChangeEvent<HTMLInputElement>, index: number) => {
//     const { value } = e.target;

//     setInputValue((prev) => prev + value);
//     if (value && index < 4) {
//       inputsRef.current[index + 1]?.focus();
//     }
//   };

//   return (
//     <form className="eachRowContainer">
//       {Array.from({ length: 5 }).map((_, index) => {
//         let backgroundColor;
//         let charColor;
//         const charResultColor: string = attempt?.result[index];
//         charColor = 'black';
//         backgroundColor = 'transparent';

//         if (charResultColor === 'G') backgroundColor = '#6aaa64';
//         if (charResultColor === 'Y') backgroundColor = '#c9b458';
//         if (charResultColor === 'X') backgroundColor = '#787c7e';
//         if (charResultColor) charColor = '#ffffff';

//         if (attempt?.result[index]) {
//           return (
//             <input
//               style={{ backgroundColor: backgroundColor, color: charColor }}
//               type="text"
//               readOnly={charResultColor ? true : false}
//               className="gameInput"
//               value={attempt?.word[index]}
//               maxLength={1}
//               // disabled={charResultColor ? true : false}
//               key={index}
//               ref={(el) => {
//                 inputsRef.current[index] = el;
//               }}
//             />
//           );
//         }
//         return (
//           <input
//             style={{ backgroundColor: backgroundColor, color: charColor }}
//             type="text"
//             readOnly={charResultColor ? true : false}
//             className="gameInput"
//             value={value[index] ? value[index] : ''}
//             maxLength={1}
//             onChange={(e) => handleChange(e, index)}
//             onKeyDown={(e) => handleKeyDown(e, index)}
//             key={index}
//             ref={(el) => {
//               inputsRef.current[index] = el;
//             }}
//           />
//         );
//       })}
//     </form>
//   );
// }

import React from 'react';
import './GameRow.css';

interface GameRowProps {
  // onWordSubmit здесь больше не нужен для работы инпутов,
  // но если он используется где-то еще, можно оставить.
  attempt?: {
    word: string;
    result: string;
  };
  value: string;
}

export default function GameRow({ attempt, value }: GameRowProps) {
  return (
    <div className="eachRowContainer">
      {Array.from({ length: 5 }).map((_, index) => {
        let backgroundColor = 'transparent';
        let charColor = 'black';
        let borderColor = '#d3d6da';

        // Получаем результат для текущей ячейки, если попытка уже отправлена
        const charResultColor = attempt?.result[index];

        if (charResultColor === 'G') backgroundColor = '#6aaa64';
        else if (charResultColor === 'Y') backgroundColor = '#c9b458';
        else if (charResultColor === 'X') backgroundColor = '#787c7e';

        if (charResultColor) {
          charColor = '#ffffff';
          borderColor = 'transparent';
        }

        // Определяем, что рисовать: букву из старой попытки или из текущего ввода
        const charToShow = attempt ? attempt.word[index] : value[index] || '';

        // Если в текущей активной ячейке есть буква, подсветим рамку (как в оригинале)
        if (!attempt && value[index]) {
          borderColor = '#878a8c';
        }

        return (
          <div
            key={index}
            className="gameCell"
            style={{
              backgroundColor,
              color: charColor,
              borderColor,
            }}
          >
            {charToShow}
          </div>
        );
      })}
    </div>
  );
}
