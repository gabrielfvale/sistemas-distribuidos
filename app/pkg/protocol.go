package pkg

var RABBITMQ_URL string = "amqp://guest:guest@localhost:5672"

type WebMessage struct {
	ActuatorType string `json:"actuatorType"`
	CommandKey   string `json:"commandKey"`
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
