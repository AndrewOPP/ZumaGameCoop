import GameRow from './GameRow/GameRow';
import './Gameview.css';
import { NetworkManager, type RoomInfo } from '../../game/network';

interface GameViewProps {
  networkManager: NetworkManager;
  roomData: RoomInfo;
}

export default function GameView({ networkManager, roomData }: GameViewProps) {
  const {
    gameState: { playerAttempts, scores },
    currentPlayerID,
  } = roomData;

  const renderRows = () => {
    const totalRows = 6;
    const attempts = playerAttempts[currentPlayerID] || [];
    return Array.from({ length: totalRows }).map((_, index) => {
      return <GameRow key={index} onWordSubmit={HandleOnSubmit} attempt={attempts[index]} />;
    });
  };

  const HandleOnSubmit = (word: string) => {
    networkManager.sendCommand(roomData.currentPlayerID, 'check_word', { word: word.toLocaleUpperCase() });
  };

  return (
    <>
      <p>Конец игры через: {roomData.gameState.timeRemaining}</p>
      <p>Score: {scores[currentPlayerID]}</p>
      <div className="rowsContainer">{renderRows()}</div>
    </>
  );
}
