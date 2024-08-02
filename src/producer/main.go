package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type Message struct {
	Date time.Time `json:"date"`
	ID   string    `json:"id"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		failOnError(err, "Failed loading settings from .env")
	}

	RABBITMQ_URL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(RABBITMQ_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()

		msg := Message{
			Date: time.Now(), ID: id,
		}

		body, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Failed to create message", http.StatusInternalServerError)
			return
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		if err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		fmt.Println(string(body))

		fmt.Fprintf(w, "Message sent: %s", body)
	})

	log.Printf("API is running on port 5001...")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
