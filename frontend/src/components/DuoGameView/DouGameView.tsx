import React, { type Dispatch, type SetStateAction } from "react";

type CurrentView = "MainMenu" | "Settings" | "Loading" | "SoloGame" | "DuoGame";

interface SoloGameViewProps {
  setCurrentView: Dispatch<SetStateAction<CurrentView>>;
}

export default function DuoGameView({ setCurrentView }: SoloGameViewProps) {
  return (
    <div>
      SoloGameView
      <p>Welcome to your DUO game! Be happy!</p>
      <button
        onClick={() => {
          setCurrentView("MainMenu");
        }}
      >
        Back
      </button>
    </div>
  );
}
