package pkg

import (
	"net"
	"time"
)

var RABBITMQ_URL string = "amqp://guest:guest@localhost:5672"

type WebMessage struct {
	ActuatorType string `json:"actuatorType"`
	CommandKey   string `json:"commandKey"`
}

type SensorMessage struct {
	Sensor    string    `json:"sensor"`
	Value     int32     `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Environment struct {
	Temperature int32
	Luminosity  int32
	Smoke       int32
}

type EnviromentMessageData struct {
	Field string
	Value int32
}

type EnvironmentMessage struct {
	Type string
	Data EnviromentMessageData
}

type Actuator struct {
	Name            string
	Health          string
	Status          bool
	EnvironmentConn net.Conn
}

type SensorInterface interface {
	Connect() error
	Consume() error
	SetEnvironment() error
}

type Sensor struct {
	Name            string
	EnvironmentConn net.Conn
}
