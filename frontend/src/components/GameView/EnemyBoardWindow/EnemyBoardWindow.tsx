import React from 'react';
import type { RoomInfo } from '../../../game/network';
import './EnemyBoardWindow.css';
interface EnemyBoardWindowProps {
  roomData: RoomInfo;
}

export default function EnemyBoardWindow({ roomData }: EnemyBoardWindowProps) {
  const { gameState, currentPlayerID } = roomData;
  const { playerAttempts, players } = gameState;

  const enemyPlayerID = Object.keys(players).filter((playerID) => playerID !== currentPlayerID)[0];
  const { score, nickname, id } = players[enemyPlayerID];

  if (id) {
    return (
      <div className="enemyDataContainer ">
        <p>Score: {score}</p>
        <p>Nickname: {nickname}</p>
        <p>Round: </p>
      </div>
    );
  }
}
