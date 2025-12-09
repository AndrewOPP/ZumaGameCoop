
// import { Renderer } from './render.js';
// import { NetworkManager } from './network.js';

// interface GameState {
// 	TestCoordinate: number,
// 	CurrentBall:    Ball
// }

// interface Ball {
//     Color:  string 
// }

// const WIDTH = 800;
// const HEIGHT = 600;

// let currentGameState = null;
// // 1. Инициализация отрисовщика
// // const renderer = new Renderer("gameCanvas", WIDTH, HEIGHT);

// // 2. Инициализация сети
// // Передаем метод renderer.draw как коллбэк для обработки полученных состояний
// const network = new NetworkManager((gameState: GameState) => {
//     // В идеале здесь можно было бы провести некую логику обновления (например, интерполяцию)
//     // но для начала просто вызываем отрисовку
//     currentGameState = gameState;
//     renderer.draw(gameState);
// });


// document.getElementById('ball-shoot').addEventListener('click', () => {
//     let color = "blue"
//     if (currentGameState.CurrentBall.color === "blue")  {
//         color = "red"
//     }
//     // В данных выстрела может быть позиция курсора или тип оружия
//     network.sendCommand("CHANGE_COLOR", {"start": 78}, {"color": color});
// });
