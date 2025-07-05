package wrapper

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqOpen interface {
	Dial(url string) (*amqp.Connection, error)
}

type RabbitMqConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
}

type RabbitMqChannel interface {
	Close() error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	Qos(prefetchCount int, prefetchSize int, global bool) error
}
