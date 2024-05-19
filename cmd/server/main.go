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
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()

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

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		"game_logs.*",
		pubsub.Durable,
	)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			log.Println("Sending a pause message")
			err = pubsub.PublishJSON(
				connCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
		case "resume":
			log.Println("Sending a resume message")
			err = pubsub.PublishJSON(
				connCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
		case "quit":
			log.Println("Shutting down server")
			return
		default:
			log.Println("I don't understand the command")
		}
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}
