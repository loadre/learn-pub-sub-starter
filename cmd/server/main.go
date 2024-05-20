package main

import (
	"fmt"
	"log"

	"github.com/loadre/learn-pub-sub-starter/internal/gamelogic"
	"github.com/loadre/learn-pub-sub-starter/internal/pubsub"
	"github.com/loadre/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const url = "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()

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

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("Could not subscrube to pause: %v:", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Sending a pause message")
			err = pubsub.PublishJSON(
				connCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
		case "resume":
			fmt.Println("Sending a resume message")
			err = pubsub.PublishJSON(
				connCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
		case "quit":
			fmt.Println("Shutting down server")
			return
		default:
			fmt.Println("I don't understand the command")
		}
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}
