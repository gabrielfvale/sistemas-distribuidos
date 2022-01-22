import { useEffect, useState } from "react";
import { Box, Button, Chip } from "@mui/material";

const socket = new WebSocket("ws://localhost:8000/ws");

function App() {
  const [isConnected, setIsConnected] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    socket.onopen = () => {
      setMessage("Connected");
      setIsConnected(true);
    };

    socket.onmessage = (e) => {
      setMessage("Get message from server: " + e.data);
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <Box>
      <Chip
        label={isConnected ? "CONNECTED" : "NOT CONNECTED"}
        color={isConnected ? "success" : "error"}
      />
      <p>{message}</p>
    </Box>
  );
}

export default App;
