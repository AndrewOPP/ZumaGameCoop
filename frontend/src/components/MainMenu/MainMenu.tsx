import { useState, type JSX, useCallback, useEffect, useRef } from 'react';
import { NetworkManager, type RoomInfo } from '../../game/network';
import './MainMenu.css';
import Loading from '../Loading/Loading';
import LobbyView from '../LobbyView/LobbyView';
import MusicButton from './MusicButton/MusicButton';
import CreateRoomModal from './CreateRoomModal/CreateRoomModal';
import JoinRoomModal from './JoinRoomModal.tsx/JoinRoomModal';
import { useGameConnection } from '../../hooks/useGameConnection';

export default function MainMenu() {
  const [openModalCreate, setOpenModalCreate] = useState<boolean>(false);
  const [openJoinModal, setOpenJoinModal] = useState<boolean>(false);
  const toggleModalCreate = () => setOpenModalCreate((prev) => !prev);
  const toggleModalJoin = () => setOpenJoinModal((prev) => !prev);

  const { currentView, roomData, networkManager, nickname, setNickname, connect, hasStoredSession, setCurrentView } =
    useGameConnection('ZumaPlayer');

  const handleCreateRoom = (roomName: string) => {
    connect('create', null, roomName);
  };

  const handleJoinRoom = (roomID: string) => {
    connect('join', roomID);
  };

  const handleReconnect = () => {
    connect('reconnect', null);
  };

  const renderView = (view: typeof currentView): JSX.Element | null => {
    switch (view) {
      case 'MainMenu':
        return null;
      case 'Loading':
        return <Loading />;
      case 'Lobby':
        // Если перешли в Lobby, но данных нет (хотя не должно быть), показываем Loading
        if (roomData && networkManager) {
          return <LobbyView roomData={roomData} networkManager={networkManager} setCurrentView={setCurrentView} />;
        }
        return <Loading />;
      default:
        return <p />;
    }
  };

  return (
    <>
      {renderView(currentView)}

      {currentView === 'MainMenu' && (
        <div className="mainMenuDiv">
          {hasStoredSession() ? (
            <button onClick={handleReconnect} style={{ backgroundColor: '#4CAF50', color: 'white' }}>
              ПЕРЕПОДКЛЮЧИТЬСЯ К СЕССИИ
            </button>
          ) : (
            <>
              <h3>Enter your Nickname:</h3>
              <input
                type="text"
                value={nickname}
                onChange={(e) => setNickname(e.target.value)}
                placeholder="Nickname"
              />
              <button onClick={toggleModalCreate} disabled={!nickname.trim()}>
                Create Room
              </button>
              <button onClick={toggleModalJoin}>Join Room</button>
            </>
          )}
          <MusicButton />

          <CreateRoomModal open={openModalCreate} onClose={toggleModalCreate} onCreateRoom={handleCreateRoom} />
          <JoinRoomModal open={openJoinModal} onClose={toggleModalJoin} connectPlayer={handleJoinRoom} />
        </div>
      )}
    </>
  );
}
