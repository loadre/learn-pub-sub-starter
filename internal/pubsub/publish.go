package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON(ch *amqp.Channel, exchange, key string, val any) error {
	dat, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("PublishJSON: %w", err)
	}

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{ContentType: "application/json", Body: dat},
	)

	return nil
}
