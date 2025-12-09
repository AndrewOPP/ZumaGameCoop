import React, { useRef, useEffect } from "react";
// Импортируйте ваши классы из папки src/game/
import { Renderer } from "../game/render";
import { NetworkManager } from "../game/network"; // Предполагаем, что NetworkManager есть

// --- Исправленные интерфейсы ---

interface Ball {
  // В классе Renderer используется r (радиус) и color (цвет)
  r: number;
  color: string;
}

interface GameState {
  TestCoordinate: number;
  CurrentBall: Ball;
  // Добавьте другие поля состояния, которые вы ожидаете от сети (например, Path, Chain)
}

// --- Компонент GameCanvas ---

const GameCanvas: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const WIDTH = 800;
  const HEIGHT = 600;

  useEffect(() => {
    // Проверка, что элемент Canvas примонтирован
    if (canvasRef.current) {
      const canvasElement = canvasRef.current;

      // 1. Инициализация отрисовщика
      // ПЕРЕДАЕМ САМ ЭЛЕМЕНТ CANVAS, а не ID.
      const renderer = new Renderer(canvasElement, WIDTH, HEIGHT);

      // 2. Инициализация сети
      // Передаем метод renderer.draw как коллбэк для обработки полученных состояний
      // Типизация коллбэка NetworkManager
      const network = new NetworkManager((gameState: GameState) => {
        // Здесь можно было бы сделать логику:
        // currentGameState = gameState;
        // Вместо сохранения в замыкании, сразу отрисовываем
        renderer.draw(gameState);
      });

      // Логика очистки (cleanup function)
      return () => {
        // Если NetworkManager имеет метод для закрытия соединения
        // network.closeConnection();
      };
    }
  }, []); // Пустой массив: запускается один раз при монтировании

  return (
    // Элемент Canvas, который React "отдаст" вашему JS-коду через ref
    <>
      <canvas ref={canvasRef} width={WIDTH} height={HEIGHT} />
      <p>Im looooxx</p>
    </>
  );
};

export default GameCanvas;
