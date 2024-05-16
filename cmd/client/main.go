package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	/// Back in the cmd/client package, use these parameters to call DeclareAndBind:
	///
	///    exchange: peril_direct (this is a constant in the internal/routing package)
	///    queueName: pause.username where username is the user's input. The pause section of the name is the routing key constant in the internal/routing package.
	///    routingKey: pause (this is a constant in the internal/routing package)
	///    simpleQueueType: transient
	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Fatalf("DeclareAndBind: %e\n", err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}
