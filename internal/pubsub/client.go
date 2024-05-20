package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueType int

const (
	_ = iota
	Durable
	Transient
)

/// Create a new .Channel() on the connection.
/// Declare a new queue using .QueueDeclare():
///
///     The durable parameter should only be true if simpleQueueType is durable.
///     The autoDelete parameter should be true if simpleQueueType is transient.
///     The exclusive parameter should be true if simpleQueueType is transient.
///     The noWait parameter should be false.
///     The args parameter should be nil.

/// Bind the queue to the exchange using .QueueBind().
/// Return the channel and queue.

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType QueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		simpleQueueType == Durable,
		simpleQueueType == Transient,
		simpleQueueType == Transient,
		false,
		nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %w", err)
	}

	err = ch.QueueBind(q.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %w", err)
	}

	return ch, q, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType QueueType,
	handler func(T),
) error {
	mqCh, q, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType)
	if err != nil {
		return fmt.Errorf("SubscribeJSON: %w", err)
	}

	// NOTE: empty consumer name s.t. it auto-generates one for us
	deliveryCh, err := mqCh.Consume(
		q.Name,
		"",
		false, false, false, false, nil)

	go func() {
		for delivery := range deliveryCh {
			var dat T
			// TODO: do we need to handle error here?
			// NOTE: remember to pass struct as a FUCKING POINTER
			_ = json.Unmarshal(delivery.Body, &dat)
			handler(dat)
			// TODO: do we need to handle error here?
			_ = delivery.Ack(false)
		}
	}()

	return nil
}
