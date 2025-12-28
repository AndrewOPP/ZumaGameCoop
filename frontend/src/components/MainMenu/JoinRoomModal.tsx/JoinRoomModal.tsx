import { useState } from 'react';
import { Box, Button, Modal, Typography, TextField } from '@mui/material';
const ModalStyle = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  border: '2px solid #000',
  boxShadow: 24,
  p: 4,
  display: 'flex',
  flexDirection: 'column',
  gap: '20px',
};

interface JoinRoomModalProps {
  open: boolean;
  onClose: () => void;
  connectPlayer: (roomID: string) => void;
}

export default function JoinRoomModal({ open, onClose, connectPlayer }: JoinRoomModalProps) {
  const [roomID, setroomID] = useState('');

  const handleConnect = () => {
    if (roomID.trim()) {
      connectPlayer(roomID.trim());
      setroomID('');
      onClose();
    }
  };

  return (
    <Modal open={open} onClose={onClose} aria-labelledby="join-room-title" aria-describedby="join-room-title">
      <Box sx={ModalStyle}>
        <Typography id="join-room-title" variant="h5" component="h2">
          Присоединиться к комнате
        </Typography>

        <TextField
          label="ID Комнаты"
          variant="outlined"
          fullWidth
          value={roomID}
          onChange={(event) => setroomID(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === 'Enter') handleConnect();
          }}
        />

        <Box sx={{ display: 'flex', justifyContent: 'space-between', gap: '10px' }}>
          <Button variant="outlined" onClick={onClose} fullWidth>
            Отмена
          </Button>
          <Button
            variant="contained"
            onClick={() => {
              connectPlayer(roomID);
              onClose();
            }}
            disabled={!roomID.trim()}
            fullWidth
          >
            Подключиться
          </Button>
        </Box>
      </Box>
    </Modal>
  );
}
