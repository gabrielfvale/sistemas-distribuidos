package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"github.com/streadway/amqp"
)

type SmokeSensor struct {
	pkg.Sensor
}

func (ls *SmokeSensor) Publish() {
	log.Printf("Publishing smoke Sensor...")
	conn, err := amqp.Dial(pkg.RABBITMQ_URL)
	pkg.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	pkg.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue_name := "smoke"

	q, err := ch.QueueDeclare(
		queue_name, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	pkg.FailOnError(err, "Failed to declare a queue")

	// Create a 5-second ticker
	for range time.Tick(5 * time.Second) {
		enviroment := pkg.ReadEnviromentData(ls.EnvironmentConn)
		message := pkg.SensorMessage{
			Sensor:    "smoke",
			Value:     enviroment.Smoke,
			Timestamp: time.Now(),
		}

		body, err := json.Marshal(message)
		pkg.FailOnError(err, "Failed to marshal json")

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		pkg.FailOnError(err, "Failed to publish a message")
		log.Printf("Sent data to queue: %v", message)
	}
}

func main() {
	smoke := SmokeSensor{}
	smoke.Name = "smoke"
	smoke.EnvironmentConn = pkg.ConnectToEnviroment()
	smoke.Publish()
}
