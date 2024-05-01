package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/loadre/learn-pub-sub-starter/internal/pubsub"
	"github.com/loadre/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	url := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close()

	connCh, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connection successful.")

	// func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	err = pubsub.PublishJSON(
		connCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)

	// wait for ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}
