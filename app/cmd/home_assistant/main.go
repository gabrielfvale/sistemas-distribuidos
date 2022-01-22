package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg/sensors"
)

var upgrader = websocket.Upgrader{}

var (
	environment        pkg.Environment
	luminosity_sensor  sensors.LuminositySensor
	smoke_sensor       sensors.SmokeSensor
	temperature_sensor sensors.TemperatureSensor
)

func load_sensors() {
	// Set sensors
	luminosity_sensor.Environment = &environment
	smoke_sensor.Environment = &environment
	temperature_sensor.Environment = &environment

	go luminosity_sensor.Publish()

	go smoke_sensor.Publish()

	go temperature_sensor.Publish()
}

func main() {
	// Set environment
	environment = pkg.Environment{Temperature: 28, Luminosity: 0, Smoke: false}
	load_sensors()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			log.Printf("Error found, %v", err)
		}
		defer ws.Close()

		log.Println("Websocket Connected!")

		for {
			var message pkg.WebMessage
			err := ws.ReadJSON(&message)
			if err != nil {
				log.Printf("Error ocurred: %v", err)
				break
			}
			log.Println(message)

			if err := ws.WriteJSON(message); err != nil {
				log.Printf("error: %v", err)
			}
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}
