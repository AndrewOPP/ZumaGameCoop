import { NetworkManager, type RoomInfo } from "../../game/network";

interface LobbyProps {
  roomData: RoomInfo;
  networkManager: NetworkManager;
  setCurrentView: (view: CurrentView) => void;
}

type CurrentView = "MainMenu" | "Loading" | "Lobby";

export default function LobbyView({
  roomData,
  networkManager,
  setCurrentView,
}: LobbyProps) {
  console.log(roomData);

  return (
    <div style={{ padding: "20px", border: "1px solid #ccc" }}>
      <h2>ROOM: {roomData.roomName}</h2>
      <h4>Room id: {roomData.roomID}</h4>
      <p>
        Your Role: <strong>{roomData.role}</strong>
      </p>

      <h3>Participants:</h3>
      <ul>
        {roomData.players.map((p) => (
          <li key={p.playerID}>
            {p.nickname} ({p.role}){p.role === "host" && <span> 👑</span>}
          </li>
        ))}
      </ul>

      <button
        onClick={() => {
          // В будущем здесь можно будет отправить команду "READY"
          // networkManager.sendCommand('READY', {});
          console.log("Sending ready signal...");
        }}
      >
        Start Game
      </button>
      <button
        onClick={() => {
          localStorage.removeItem("roomID");
          localStorage.removeItem("playerID");
        }}
      >
        LEAVE ROOM
      </button>
    </div>
  );
}
