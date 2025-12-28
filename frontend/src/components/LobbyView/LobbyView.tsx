import { useState } from 'react';
import { NetworkManager, type RoomInfo } from '../../game/network';
import GameView from '../GameView/GameView';

interface LobbyProps {
  roomData: RoomInfo;
  networkManager: NetworkManager;
  setCurrentView: (view: CurrentView) => void;
}

type CurrentView = 'MainMenu' | 'Loading' | 'Lobby';

export default function LobbyView({ roomData, networkManager, setCurrentView }: LobbyProps) {
  console.log(roomData);
  // const [roundTimer, setRoundTimer] = useState<number>(roomData.gameState.timeRemaining);

  return (
    <>
      {roomData.gameState.isActive ? (
        <GameView networkManager={networkManager} roomData={roomData} />
      ) : (
        <div style={{ padding: '20px', border: '1px solid #ccc' }}>
          <h2>ROOM: {roomData.roomName}</h2>
          <h4>Room id: {roomData.roomID}</h4>
          <p>
            Your Role: <strong>{roomData.role}</strong>
          </p>

          <h3>Participants:</h3>
          <ul>
            {roomData.players.map((player) => (
              <li key={player.id}>
                <span
                  style={{
                    color: player.isReady ? '#4caf50' : 'inherit', // Зеленый если готов, иначе обычный
                    fontWeight: player.isReady ? 'bold' : 'normal',
                  }}
                >
                  {player.nickname}
                </span>
                ({player.role}){player.role === 'host' && <span> 👑</span>}
              </li>
            ))}
          </ul>
          <button
            onClick={() => {
              setCurrentView('MainMenu');
              localStorage.removeItem('roomID');
              localStorage.removeItem('playerID');
            }}
          >
            Leave room
          </button>
          <button
            onClick={() => {
              // В будущем здесь можно будет отправить команду "READY"
              networkManager.sendCommand(roomData.currentPlayerID, 'toggle_ready');
              console.log('Sending ready signal...');
            }}
          >
            Ready
          </button>
        </div>
      )}
    </>
  );
}
