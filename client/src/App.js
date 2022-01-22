import { useEffect, useState } from "react";
import { Box, Button, Chip } from "@mui/material";

import AlarmLightOutline from "mdi-material-ui/AlarmLightOutline";
import AlarmLightOffOutline from "mdi-material-ui/AlarmLightOffOutline";
import Radiator from "mdi-material-ui/Radiator";
import RadiatorOff from "mdi-material-ui/RadiatorOff";
import LightbulbOnOutline from "mdi-material-ui/LightbulbOnOutline";
import LightbulbOffOutline from "mdi-material-ui/LightbulbOffOutline";

const socket = new WebSocket("ws://localhost:8000/ws");

function App() {
  const [isConnected, setIsConnected] = useState(false);
  const [message, setMessage] = useState("");
  const [actuatorOn, setActuatorOn] = useState([false, false, false]);

  const toggleActuator = (key, index) => {
    const prev = [...actuatorOn];
    const newValue = !prev[index];

    prev[index] = newValue;
    setActuatorOn([...prev]);

    console.log({
      actuatorType: key,
      commandKey: newValue ? "TurnOn" : "TurnOff",
    });

    socket.send(
      JSON.stringify({
        actuatorType: key,
        commandKey: newValue ? "TurnOn" : "TurnOff",
      })
    );
  };

  const actuators = [
    {
      name: "Fire alarm",
      actuatorType: "fire",
      actions: ["TurnOn", "TurnOff"],
      icons: { on: AlarmLightOutline, off: AlarmLightOffOutline },
    },
    {
      name: "Heater",
      actuatorType: "heater",
      actions: ["TurnOn", "TurnOff"],
      icons: { on: Radiator, off: RadiatorOff },
    },
    {
      name: "Lamp",
      actuatorType: "lamp",
      actions: ["TurnOn", "TurnOff"],
      icons: { on: LightbulbOnOutline, off: LightbulbOffOutline },
    },
  ];

  useEffect(() => {
    socket.onopen = () => {
      setMessage("Connected");
      setIsConnected(true);
    };

    socket.onmessage = (e) => {
      console.log("message");
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
      {actuators.map((actuator, index) => (
        <>
          <Button onClick={() => toggleActuator(actuator.actuatorType, index)}>
            {actuatorOn[index] ? <actuator.icons.on /> : <actuator.icons.off />}
          </Button>
        </>
      ))}
    </Box>
  );
}

export default App;
