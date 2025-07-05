package wrapper

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqWrapper struct {
	RabbitMqConnectionWrapper
	RabbitMqChannelWrapper
}

type RabbitMqConnectionWrapper struct {
	conn *amqp.Connection
	RabbitMqConnection
}

type RabbitMqChannelWrapper struct {
	channel *amqp.Channel
	RabbitMqChannel
}

func (rw *RabbitMqWrapper) Dial(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func (co *RabbitMqConnectionWrapper) Channel() (*amqp.Channel, error) {
	return co.conn.Channel()
}

func (co *RabbitMqConnectionWrapper) Close() error {
	return co.conn.Close()
}

func (co *RabbitMqConnectionWrapper) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return co.conn.NotifyClose(receiver)
}

func (ch *RabbitMqChannelWrapper) Close() error {
	return ch.channel.Close()
}

func (ch *RabbitMqChannelWrapper) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return ch.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (ch *RabbitMqChannelWrapper) PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return ch.channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

func (ch *RabbitMqChannelWrapper) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return ch.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (ch *RabbitMqChannelWrapper) Qos(prefetchCount int, prefetchSize int, global bool) error {
	return ch.channel.Qos(prefetchCount, amqp.PreconditionFailed, global)
}
