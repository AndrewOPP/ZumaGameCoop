import React, { useEffect, useRef, useState } from "react";

export default function MusicButton() {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);

  useEffect(() => {
    if (audioRef.current) {
      // Устанавливаем громкость на 20% (0.2)
      audioRef.current.volume = 0.1;
    }
  }, []);

  const handleTogglePlay = () => {
    if (!audioRef.current) return;

    if (isPlaying) {
      // 1. Если играет -> Пауза
      audioRef.current.pause();
      setIsPlaying(false);
    } else {
      // 2. Если не играет -> Запуск
      // Браузер разрешит play(), так как это происходит по клику пользователя
      audioRef.current
        .play()
        .then(() => {
          setIsPlaying(true);
        })
        .catch((error) => {
          // Если возникла ошибка (например, файл не найден), мы ее обрабатываем
          console.error("Не удалось запустить аудио:", error);
          setIsPlaying(false); // Сбрасываем состояние
        });
    }
  };
  return (
    <div>
      {/* 
        Элемент <audio>
        - loop: чтобы песня повторялась
        - ref: для доступа через React
      */}
      <audio ref={audioRef} src="/audio/52.mp3" loop />

      {/* Кнопка запуска/паузы */}
      <button
        onClick={handleTogglePlay}
        style={{ padding: "15px 30px", fontSize: "18px", cursor: "pointer" }}
      >
        {isPlaying ? "⏸️ Выключить" : "▶️ Включить музыку"}
      </button>

      {/* 
        ВАЖНО: Добавьте ваш файл '/music/main_theme.mp3' в 
        папку public вашего React-проекта или используйте правильный путь. 
      */}
    </div>
  );
}
