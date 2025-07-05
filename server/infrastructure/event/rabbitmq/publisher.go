package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	"github.com/kurnhyalcantara/koer-tax-service/config"
	rabbitmqwrapper "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/event/rabbitmq/wrapper"
)

type Publisher struct {
	Ctx       context.Context
	Name      string
	Queue     string
	Config    *Config
	Amqpw     rabbitmqwrapper.RabbitMqOpen
	AppConfig *config.Config
}

func (p *Publisher) InitPublisher() error {

	publisherAmqp := NewConnection(p.Name, p.Queue, p.Config, p.Amqpw)

	if err := publisherAmqp.Connect(); err != nil {
		return fmt.Errorf("failed connect mq with queue %s: %v", p.Queue, err)
	}

	if err := publisherAmqp.InitChannel(); err != nil {
		return fmt.Errorf("failed init channel mq with queue %s: %v", p.Queue, err)
	}

	if err := publisherAmqp.BindQueue(); err != nil {
		return fmt.Errorf("failed bind queue %s: %v", p.Queue, err)
	}

	return nil
}

func (p *Publisher) PublishQueue(contentType string, qmsg []byte) error {
	rabbitPublisher := GetConnection(p.Name)
	if rabbitPublisher == nil {
		return errors.New("rabbitmq publisher " + p.Name + " connection not found")
	}

	if publishErr := rabbitPublisher.Publish(p.Ctx, contentType, qmsg); publishErr != nil {
		return publishErr
	}

	return nil
}

func (p *Publisher) CleanupPublisher() error {
	conn := GetConnection(p.Name)
	if conn == nil {
		return errors.New(p.Name + " connection already closed")
	}
	conn.CleanUp()

	return nil
}
