package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	pubsub "github.com/loadre/learn-pub-sub-starter/internal/pubsub"
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

	// wait for ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}
