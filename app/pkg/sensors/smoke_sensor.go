package sensors

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"github.com/streadway/amqp"
)

type SmokeSensor struct {
	pkg.Sensor
	Smoke bool
}

type SmokeSensorMessage struct {
	Value     bool
	Timestamp time.Time
}

func (ls *SmokeSensor) Publish() {
	log.Printf("Publishing smoke Sensor...")
	conn, err := amqp.Dial(pkg.RABBITMQ_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare a queue")

	// Create a 5-second ticker
	for range time.Tick(5 * time.Second) {
		body, err := json.Marshal(SmokeSensorMessage{Value: ls.Environment.Smoke, Timestamp: time.Now()})
		failOnError(err, "Failed to marshal json")

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		failOnError(err, "Failed to publish a message")
	}
}