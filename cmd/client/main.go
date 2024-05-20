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
	fmt.Println("Starting Peril client...")

	url := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()

	mqCh, q, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Fatalf("failed to declare and bind: %v\n", err)
		return
	}
	log.Printf("Queue %v declared and bound\n", q.Name)

	gs := gamelogic.NewGameState(username)
	_ = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gs))

	_ = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.Transient,
		handlerMove(gs))

	for {
		var err error
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err = gs.CommandSpawn(words)
		case "move":
			var move gamelogic.ArmyMove
			move, err = gs.CommandMove(words)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			err = pubsub.PublishJSON(
				mqCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+username,
				move)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			fmt.Println("Published move successfully")
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			// if any other command is entered,
			// print an error message and continue the loop
		}
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}
