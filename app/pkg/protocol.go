package pkg

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
	Environment *Environment
}

type SensorInterface interface {
	Connect()
	SetEnvironment()
}

type Sensor struct {
	Name        string
	Environment *Environment
}
