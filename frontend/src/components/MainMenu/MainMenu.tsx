import { useState, type JSX } from "react";
import Loading from "../Loading/Loading";
import SoloGameView from "../SoloGameView/SoloGameView";
import DouGameView from "../DuoGameView/DouGameView";
import MusicButton from "./MusicButton/MusicButton";

export default function MainMenu() {
  type CurrentView =
    | "MainMenu"
    | "Settings"
    | "Loading"
    | "SoloGame"
    | "DuoGame";
  const [currentView, setCurrentView] = useState<CurrentView>("MainMenu");

  const renderView = (view: CurrentView): JSX.Element | null => {
    switch (view) {
      case "MainMenu":
        return null;
      case "Loading":
        return <Loading />;
      case "SoloGame":
        return <SoloGameView setCurrentView={setCurrentView} />;
      case "DuoGame":
        return <DouGameView setCurrentView={setCurrentView} />;
      default:
        return <p />; // Обязательный дефолт
    }
  };

  return (
    <>
      {renderView(currentView)}
      {currentView === "MainMenu" && (
        <>
          <button onClick={() => setCurrentView("SoloGame")}>Solo Game</button>
          <button onClick={() => setCurrentView("DuoGame")}>Duo Game</button>
          <button onClick={() => setCurrentView("Settings")}>Settings</button>
          <MusicButton />
        </>
      )}
    </>
  );
}
