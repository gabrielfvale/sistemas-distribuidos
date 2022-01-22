package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/gabrielfvale/ti0151-sistemas/app/grpc/proto"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg/actuators"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg/sensors"
)

var upgrader = websocket.Upgrader{}

var (
	environment        pkg.Environment
	luminosity_sensor  sensors.LuminositySensor
	smoke_sensor       sensors.SmokeSensor
	temperature_sensor sensors.TemperatureSensor

	fire_actuator   actuators.FireAlarmServer
	heater_actuator actuators.HeaterServer
	lamp_actuator   actuators.LampServer
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

func load_actuators(ports map[string]int) {
	fire_actuator.Environment = &environment
	heater_actuator.Environment = &environment
	lamp_actuator.Environment = &environment

	go fire_actuator.Listen(ports["fire"])
	go heater_actuator.Listen(ports["heater"])
	go lamp_actuator.Listen(ports["lamp"])

}

func load_queues() {

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func processQueue(name string, ws *websocket.Conn) {
	conn, _ := amqp.Dial(pkg.RABBITMQ_URL)
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	msgs, _ := ch.Consume(name, "home_assistant", true, false, false, false, nil)
	go func() {
		for d := range msgs {
			ws.WriteJSON([]byte(d.Body))
		}
	}()
}

func main() {
	// Set environment
	// queues := []string{"luminosity", "smoke", "temperature"}

	actuator_ports := make(map[string]int)
	actuator_ports["fire"] = 8001
	actuator_ports["heater"] = 8002
	actuator_ports["lamp"] = 8003

	environment = pkg.Environment{Temperature: 28, Luminosity: 0, Smoke: false}

	load_sensors()
	load_actuators(actuator_ports)

	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", actuator_ports["heater"]), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewActuatorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.GetAvailableCommands(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

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

		// forever := make(chan bool)
		// for _, queue := range queues {
		// 	go processQueue(queue, ws)
		// }

		for {
			var message pkg.WebMessage
			err = ws.ReadJSON(&message)
			if err != nil {
				log.Printf("Error ocurred: %v", err)
				break
			}
			log.Println(message)

			var res *pb.IssueCommandResponse

			switch message.ActuatorType {
			case "fire":
				res, _ = fire_actuator.IssueCommand(context.Background(), &pb.IssueCommandRequest{Key: message.CommandKey})
			case "heater":
				res, _ = heater_actuator.IssueCommand(context.Background(), &pb.IssueCommandRequest{Key: message.CommandKey})
			case "lamp":
				res, _ = lamp_actuator.IssueCommand(context.Background(), &pb.IssueCommandRequest{Key: message.CommandKey})
			}

			if err := ws.WriteJSON(res); err != nil {
				log.Printf("error: %v", err)
			}
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}
