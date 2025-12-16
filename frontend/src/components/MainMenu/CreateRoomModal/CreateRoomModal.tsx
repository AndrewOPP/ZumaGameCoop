import React, { useState } from "react";
import { Box, Button, Modal, Typography, TextField } from "@mui/material";

// --- СТИЛИ (Вынесены в константу, использующую MUI's sx) ---
const ModalStyle = {
  position: "absolute", // Явное указание типа для TypeScript
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 400,
  bgcolor: "background.paper",
  border: "2px solid #000",
  boxShadow: 24,
  p: 4,
  display: "flex",
  flexDirection: "column",
  gap: "20px",
};

// Интерфейс для пропсов
interface CreateRoomModalProps {
  open: boolean;
  onClose: () => void;
  onCreateRoom: (roomName: string) => void;
}

export default function CreateRoomModal({
  open,
  onClose,
  onCreateRoom,
}: CreateRoomModalProps) {
  const [roomName, setRoomName] = useState("");

  const handleCreate = () => {
    if (roomName.trim()) {
      onCreateRoom(roomName.trim());
      setRoomName(""); // Очищаем поле
      onClose(); // Закрываем модальное окно
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby="create-room-title"
      aria-describedby="create-room-description"
    >
      <Box sx={ModalStyle}>
        <Typography id="create-room-title" variant="h5" component="h2">
          Создать новую комнату
        </Typography>

        <TextField
          label="Название комнаты"
          variant="outlined"
          fullWidth
          value={roomName}
          onChange={(event) => setRoomName(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") handleCreate();
          }}
        />

        <Box
          sx={{ display: "flex", justifyContent: "space-between", gap: "10px" }}
        >
          <Button variant="outlined" onClick={onClose} fullWidth>
            Отмена
          </Button>
          <Button
            variant="contained"
            onClick={handleCreate}
            disabled={!roomName.trim()}
            fullWidth
          >
            Создать
          </Button>
        </Box>
      </Box>
    </Modal>
  );
}
