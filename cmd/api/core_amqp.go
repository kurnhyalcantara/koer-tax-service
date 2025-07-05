package main

import (
	"context"
	"log"
	"strconv"

	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/event/rabbitmq"
	rabbitmqwrapper "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/event/rabbitmq/wrapper"
)

func InitPublisherQueue(ctx context.Context) *rabbitmq.Publisher {
	reconnectDelay, _ := strconv.Atoi(appConfig.AmqpReconnectDelay)

	publisher := &rabbitmq.Publisher{
		Ctx:   ctx,
		Name:  appConfig.AmqpRegistrationOnlineQueue,
		Queue: appConfig.AmqpRegistrationOnlineQueue,
		Config: &rabbitmq.Config{
			User:           appConfig.AmqpUser,
			Password:       appConfig.AmqpPassword,
			Host:           appConfig.AmqpHost,
			ReconnectDelay: reconnectDelay,
			QueueConfig: &rabbitmq.QueueConfig{
				Durable:       true,
				AutoDelete:    false,
				Exclusive:     false,
				NoWait:        false,
				PrefetchCount: 1,
				PrefetchSize:  0,
			},
		},
		Amqpw:     &rabbitmqwrapper.RabbitMqWrapper{},
		AppConfig: appConfig,
	}

	if initGetErr := publisher.InitPublisher(); initGetErr != nil {
		log.Fatalf("failed to initialize publisher queue: %v", initGetErr)
	}

	return publisher
}
