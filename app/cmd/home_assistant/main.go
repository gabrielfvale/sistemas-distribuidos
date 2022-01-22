package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			log.Printf("Error found, %v", err)
		}
		defer ws.Close()

		log.Println("Connected!")

		for {
			var message Message
			err := ws.ReadJSON(&message)
			if err != nil {
				log.Printf("Error ocurred: %v", err)
				break
			}
			log.Printf(message.Message)

			if err := ws.WriteJSON(message); err != nil {
				log.Printf("error: %v", err)
			}
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}
