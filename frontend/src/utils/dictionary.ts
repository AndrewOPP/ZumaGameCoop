import rawWords from '../words.txt?raw'; // Если используете Vite

// Превращаем строку в Set для мгновенного поиска O(1)
const wordsArray = rawWords.split('\n').map((w) => w.trim().toUpperCase());
export const dictionary = new Set(wordsArray);

export const isValidWord = (word: string): boolean => {
  return dictionary.has(word.toUpperCase());
};
