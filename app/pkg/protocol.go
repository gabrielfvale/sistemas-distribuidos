package pkg

import "google.golang.org/grpc"

var RABBITMQ_URL string = "amqp://guest:guest@localhost:5672"

type WebMessage struct {
	Message string `json:"message"`
}

type Environment struct {
	Temperature uint
	Luminosity  uint32
	Smoke       bool
}

type Actuator struct {
	Name        string
	Health      string
	Status      bool
	Server      *grpc.Server
	Environment *Environment
}

type SensorInterface interface {
	Connect() error
	Consume() error
	SetEnvironment() error
}

type Sensor struct {
	Name        string
	Environment *Environment
}
