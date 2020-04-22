package amqp

import "github.com/streadway/amqp"

type amqpEventListener struct {
	connection *amqp.Connection
}
