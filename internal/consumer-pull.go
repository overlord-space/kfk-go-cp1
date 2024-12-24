package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

const pollTimeoutMs = 100

func CreatePullCustomer(daemonData *DaemonSharedData) error {
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers":  BootstrapServers,
		"group.id":           "pull-example-group",
		"auto.offset.reset":  "earliest",
		"session.timeout.ms": 6000,
	}

	consumer, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("[pull] Не удалось создать консьюмера: %s\n", err))
	}

	log.Printf("[pull] Консьюмер создан %v\n", consumer)

	err = consumer.SubscribeTopics([]string{TopicName}, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("[pull] Не удалось подписаться на топик: %s\n", err))
	}

	isRunning := true

	for isRunning {
		select {
		case <-daemonData.ExistSignalChannel:
			isRunning = false
		default:
		}

		rawEvent := consumer.Poll(pollTimeoutMs)
		if rawEvent == nil {
			continue
		}

		switch event := rawEvent.(type) {
		case *kafka.Message:
			log.Printf("[pull] Получено сообщение: %v\n", string(event.Value))

			messageValue := SomeMessage{}
			err := json.Unmarshal(event.Value, &messageValue)
			if err != nil {
				log.Printf("[pull] Не удалось десериализовать сообщение: %s\n", err)
			} else {
				log.Printf("[pull][json] Получено сообщение: %v\n", messageValue)
			}
		default:
		}
	}

	log.Printf("[pull] Завершение работы консьюмера\n")
	err = consumer.Close()
	if err != nil {
		log.Printf("Не удалось завершить работу консьюмера: %s\n", err)
	}

	daemonData.WaitGroup.Done()
	return nil
}
