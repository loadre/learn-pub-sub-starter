package pubsub

import (
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
	simpleQueueType int, // an enum to represent "durable" or "transient"
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
		return nil, amqp.Queue{}, fmt.Errorf("QueueDeclare: %w", err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("QueueBind: %w", err)
	}

	return ch, q, nil
}
