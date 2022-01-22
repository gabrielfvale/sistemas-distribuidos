package impl

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
