package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	logger *slog.Logger
)

type Message struct {
	Date time.Time `json:"date"`
	ID   string    `json:"id"`
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", msg, err.Error()))
		os.Exit(1)
	}
}

func main() {
	logLevel := parseLogLevel("DEBUG")

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger = slog.New(logHandler)
	slog.SetDefault(logger)

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

		logger.Debug(fmt.Sprintf("Message sent: %s", string(body)))
	})

	logger.Info("API is running on port 5001...")

	if err := http.ListenAndServe(":5001", nil); err != nil {
		failOnError(err, "Failed to start HTTP server")
	}
}

func parseLogLevel(logLevel string) slog.Level {
	logLevel = strings.ToUpper(logLevel)
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
