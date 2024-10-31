package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	logger *slog.Logger
)

func failOnError(err error, msg string) {
	if err != nil {
		logger.Error("%s: %s", msg, err)
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logger.Debug(fmt.Sprintf("Received a message: %s", string(d.Body)))
		}
	}()

	logger.Info("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
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
