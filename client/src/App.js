import { useEffect, useState } from "react";
import { Box, Button, Chip, Stack } from "@mui/material";

import AlarmLightOutline from "mdi-material-ui/AlarmLightOutline";
import AlarmLightOffOutline from "mdi-material-ui/AlarmLightOffOutline";
import Radiator from "mdi-material-ui/Radiator";
import RadiatorOff from "mdi-material-ui/RadiatorOff";
import LightbulbOnOutline from "mdi-material-ui/LightbulbOnOutline";
import LightbulbOffOutline from "mdi-material-ui/LightbulbOffOutline";
import { Thermometer, Home, Plus, Minus } from "mdi-material-ui";

const socket = new WebSocket("ws://localhost:8000/ws");

function App() {
  const [isConnected, setIsConnected] = useState(false);
  const [actuatorOn, setActuatorOn] = useState([false, false, false]);
  const [temperature, setTemperature] = useState(0);

  const sendAction = (key, action) => {
    socket.send(
      JSON.stringify({
        actuatorType: key,
        commandKey: action,
      })
    );
  };

  const toggleActuator = (key, index) => {
    const prev = [...actuatorOn];
    const newValue = !prev[index];

    prev[index] = newValue;
    setActuatorOn([...prev]);

    console.log({
      actuatorType: key,
      commandKey: newValue ? "TurnOn" : "TurnOff",
    });
    sendAction(key, newValue ? "TurnOn" : "TurnOff");
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
      actions: ["TurnOn", "TurnOff", "RaiseTemp", "LowerTemp"],
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
      // setMessage("Connected");
      setIsConnected(true);
    };

    socket.onmessage = (e) => {
      const data = JSON.parse(e.data);

      if (data?.sensor) {
        setTemperature(data.value);
      }
      console.log(data);
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        width: "100%",
        height: "100vh",
        backgroundColor: "#ccc",
      }}
    >
      <Stack
        direction="row"
        spacing={2}
        sx={{ display: "flex", alignItems: "center" }}
      >
        <Home />

        <h2>Welcome to Home Assistant</h2>
        <Chip
          label={isConnected ? "CONNECTED" : "NOT CONNECTED"}
          color={isConnected ? "success" : "error"}
        />
      </Stack>

      <Box sx={{ display: "flex", alignItems: "center" }}>
        <Thermometer />
        <p>{temperature} °C</p>
      </Box>

      <Stack direction="row" spacing={8}>
        {actuators.map((actuator, index) => (
          <Stack spacing={2}>
            <Button
              variant="contained"
              startIcon={
                actuatorOn[index] ? (
                  <actuator.icons.on />
                ) : (
                  <actuator.icons.off />
                )
              }
              onClick={() => toggleActuator(actuator.actuatorType, index)}
            >
              {actuator.name}
            </Button>
            {actuator.actuatorType === "heater" && (
              <>
                <Button
                  variant="contained"
                  startIcon={<Plus />}
                  onClick={() => sendAction(actuator.actuatorType, "RaiseTemp")}
                >
                  Increase
                </Button>
                <Button
                  variant="contained"
                  startIcon={<Minus />}
                  onClick={() => sendAction(actuator.actuatorType, "LowerTemp")}
                >
                  Decrease
                </Button>
              </>
            )}
          </Stack>
        ))}
      </Stack>
    </Box>
  );
}

export default App;
