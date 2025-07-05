package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/event/rabbitmq/wrapper"
	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate mockery --name=RabbitMqCore --output=../../tests/mocks --structname=MockRabbitMqCore

type RabbitMqCore interface {
	GetName() string
	CheckConnection() chan error
	Connect() error
	InitChannel() error
	BindQueue() error
	CloseChannel() error
	Reconnect() error
	Consume() (<-chan amqp.Delivery, error)
	Publish(ctx context.Context, contentType string, qmsg []byte) error
	CleanUp() error
}

type Connection struct {
	name    string
	queue   string
	config  *Config
	err     chan error
	amqpw   wrapper.RabbitMqOpen
	channel wrapper.RabbitMqChannel
	conn    wrapper.RabbitMqConnection
}

type QueueConfig struct {
	// #region Queue Declare Config
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	// #endregion

	// #region Qos Config
	PrefetchCount int
	PrefetchSize  int
	Global        bool
	// #endregion
}

type Config struct {
	User           string
	Password       string
	Host           string
	ReconnectDelay int
	*QueueConfig
}

var (
	connectionPool = make(map[string]*Connection)
)

func NewConnection(name string, queue string, config *Config, amqpw wrapper.RabbitMqOpen) RabbitMqCore {
	if c, ok := connectionPool[name]; ok {
		return c
	}

	c := &Connection{
		name:   name,
		queue:  queue,
		config: config,
		err:    make(chan error),
		amqpw:  amqpw,
	}
	connectionPool[name] = c

	return c
}

func GetConnection(name string) *Connection {
	return connectionPool[name]
}

func (c *Connection) GetName() string {
	return c.name
}

func (c *Connection) CheckConnection() chan error {
	return c.err
}

func (c *Connection) Connect() error {
	var err error
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s/",
		c.config.User,
		c.config.Password,
		c.config.Host)
	c.conn, err = c.amqpw.Dial(rabbitUrl)
	if err != nil {
		return err
	}

	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()

	return nil
}

func (c *Connection) InitChannel() error {
	var channelMqErr error
	c.channel, channelMqErr = c.conn.Channel()
	if channelMqErr != nil {
		return channelMqErr
	}

	return nil
}

func (c *Connection) CloseChannel() error {
	err := c.channel.Close()
	if err != nil {
		return errors.New("channel already closed")
	}

	return nil
}

func (c *Connection) BindQueue() error {
	if _, err := c.channel.QueueDeclare(c.queue, c.config.Durable, c.config.AutoDelete, c.config.Exclusive, c.config.NoWait, c.config.Args); err != nil {
		return err
	}

	if err := c.channel.Qos(c.config.PrefetchCount, c.config.PrefetchSize, c.config.Global); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Reconnect() error {
	if cErr := c.Connect(); cErr != nil {
		return cErr
	}
	if icErr := c.InitChannel(); icErr != nil {
		return icErr
	}
	if bqErr := c.BindQueue(); bqErr != nil {
		return bqErr
	}
	return nil
}

func (c *Connection) Consume() (<-chan amqp.Delivery, error) {
	if c.channel == nil {
		return nil, errors.New("consumer channel is nil, please init consumer first")
	}
	return c.channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c *Connection) Publish(ctx context.Context, contentType string, qmsg []byte) error {
	if c.channel == nil {
		return errors.New("producer channel is nil, please init producer first")
	}

	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			if reconErr := c.Reconnect(); reconErr != nil {
				return reconErr
			}
		}
	default:
	}

	proErr := c.channel.PublishWithContext(ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         qmsg,
		})

	if proErr != nil {
		return proErr
	}

	return nil
}

func (c *Connection) CleanUp() error {
	if chErr := c.CloseChannel(); chErr != nil {
		return chErr
	}

	if coErr := c.conn.Close(); coErr != nil {
		return errors.New("connection already closed")
	}
	delete(connectionPool, c.name)
	return nil
}
