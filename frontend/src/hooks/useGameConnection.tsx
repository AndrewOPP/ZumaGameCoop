import { useState, useCallback, useEffect } from 'react';
import { NetworkManager, type RoomInfo } from '../game/network';

type ConnectionAction = 'create' | 'join' | 'reconnect';
type ViewStatus = 'MainMenu' | 'Loading' | 'Lobby';

interface ConnectionResult {
  currentView: ViewStatus;
  networkManager: NetworkManager | null;
  nickname: string;
  roomData: RoomInfo | null;
  setNickname: React.Dispatch<React.SetStateAction<string>>;
  connect: (action: ConnectionAction, roomID: string | null, roomName?: string) => void;
  hasStoredSession: () => boolean | '' | null;
  setCurrentView: React.Dispatch<React.SetStateAction<ViewStatus>>; // Полезно для выхода из лобби
}

export const useGameConnection = (initialNickname: string): ConnectionResult => {
  const [currentView, setCurrentView] = useState<ViewStatus>('MainMenu');
  const [roomData, setRoomData] = useState<RoomInfo | null>(null);
  const [nickname, setNickname] = useState<string>(initialNickname);
  const [networkManager, setNetworkManager] = useState<NetworkManager | null>(null);

  // 1. КОЛЛБЭК ДЛЯ ОБРАБОТКИ СООБЩЕНИЙ С СЕРВЕРА
  const handleServerMessage = useCallback((data: RoomInfo & { message?: string; error?: string }) => {
    if (data.type === 'room_created' || data.type === 'room_reconnected' || data.type === 'room_joined') {
      setRoomData(data);
      setCurrentView('Lobby');
      localStorage.setItem('roomID', data.roomID);
      localStorage.setItem('playerID', data.currentPlayerID);
    } else if (data.type === 'room_updated') {
      setRoomData((prevData) => {
        if (!prevData) return null;
        return { ...prevData, players: data.players };
      });
    } else if (data.type === 'error') {
      console.error('Server Error:', data.error);
      setCurrentView('MainMenu');
      alert(`Connection Error: ${data.error}`);
      localStorage.removeItem('roomID');
      localStorage.removeItem('playerID');
    } else if (data.type === 'STATUS') {
      console.log(`Connection Status: ${data.message}`);
    }
  }, []);

  const connect = useCallback(
    (action: ConnectionAction, roomID: string | null, roomName?: string) => {
      if (!nickname.trim()) {
        alert('Please enter a nickname.');
        return;
      }

      setCurrentView('Loading');

      const nm = new NetworkManager(handleServerMessage);
      setNetworkManager(nm);

      let playerID: string | null = null;
      let targetRoomID: string | null = roomID;

      if (action === 'reconnect') {
        playerID = localStorage.getItem('playerID');
        targetRoomID = localStorage.getItem('roomID');
        if (!playerID || !targetRoomID) {
          setCurrentView('MainMenu');
          alert('Нет сохраненной сессии.');
          return;
        }
      }
      nm.connect(action, playerID, targetRoomID, nickname, roomName);
    },
    [nickname, handleServerMessage],
  );

  const hasStoredSession = useCallback(() => {
    const roomID = localStorage.getItem('roomID');
    const playerID = localStorage.getItem('playerID');
    return roomID && playerID !== 'null';
  }, []);

  return {
    currentView,
    roomData,
    networkManager,
    nickname,
    setNickname,
    connect: (action, roomID, roomName) => {
      if (!nickname.trim()) {
        alert('Please enter a nickname.');
        setCurrentView('MainMenu');
        return;
      }
      connect(action, roomID, roomName);
    },
    hasStoredSession,
    setCurrentView,
  };
};
