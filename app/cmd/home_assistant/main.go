package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/gabrielfvale/ti0151-sistemas/app/grpc/proto"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
)

var upgrader = websocket.Upgrader{}

func main() {
	// Set environment
	actuator_ports := make(map[string]int)
	actuator_ports["fire"] = 8001
	actuator_ports["heater"] = 8002
	actuator_ports["lamp"] = 8003

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

		conn, _ := amqp.Dial(pkg.RABBITMQ_URL)
		defer conn.Close()

		ch, _ := conn.Channel()
		defer ch.Close()

		// luminosity_msgs, _ := ch.Consume("luminosity", "home_assistant", true, false, false, false, nil)
		// smoke_msgs, _ := ch.Consume("smoke", "home_assistant", true, false, false, false, nil)
		temperature_msgs, _ := ch.Consume("temperature", "home_assistant", true, false, false, false, nil)
		go func() {
			for d := range temperature_msgs {
				rec := pkg.SensorMessage{}
				json.Unmarshal(d.Body, &rec)
				ws.WriteJSON(rec)
			}
		}()

		for {
			var message pkg.WebMessage
			err = ws.ReadJSON(&message)
			if err != nil {
				log.Printf("Error ocurred: %v", err)
				break
			}
			log.Println(message, actuator_ports[message.ActuatorType])

			var res *pb.IssueCommandResponse
			conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", actuator_ports[message.ActuatorType]), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewActuatorClient(conn)
			c.IssueCommand(context.Background(), &pb.IssueCommandRequest{
				Key: message.CommandKey,
			})

			if err := ws.WriteJSON(res); err != nil {
				log.Printf("error: %v", err)
			}
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}
