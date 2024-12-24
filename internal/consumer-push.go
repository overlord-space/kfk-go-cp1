package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

func CreatePushConsumer(daemonData *DaemonSharedData) error {
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers":  BootstrapServers,
		"group.id":           "push-example-group",
		"auto.offset.reset":  "earliest",
		"session.timeout.ms": 6000,
		"enable.auto.commit": false,
	}

	consumer, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("[push] Не удалось создать консьюмера: %s\n", err))
	}

	log.Printf("[push] Консьюмер создан %v\n", consumer)

	err = consumer.SubscribeTopics([]string{TopicName}, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("[push] Не удалось подписаться на топик: %s\n", err))
	}

	isRunning := true

	for isRunning {
		select {
		case <-daemonData.ExistSignalChannel:
			isRunning = false
		case message := <-daemonData.PushChannel:
			log.Printf("[push] Получено сообщение: %v\n", string(message.Value))

			messageValue := SomeMessage{}
			err := json.Unmarshal(message.Value, &messageValue)
			if err != nil {
				log.Printf("[push] Не удалось десериализовать сообщение: %s\n", err)
			} else {
				log.Printf("[push][json] Получено сообщение: %v\n", messageValue)

				_, err := consumer.CommitMessage(message)
				if err != nil {
					return errors.New(fmt.Sprintf("[push] Не удалось подтвердить получение сообщения: %s\n", err))
				}
			}
		}
	}

	log.Print("[push] Завершение работы консьюмера")
	daemonData.WaitGroup.Done()
	return nil
}
