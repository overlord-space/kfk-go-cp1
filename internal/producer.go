package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"time"
)

func CreateProducer(daemonData *DaemonSharedData) error {
	//goland:noinspection SpellCheckingInspection
	producerConfig := kafka.ConfigMap{
		"bootstrap.servers": BootstrapServers,
		"acks":              "all",
		"retries":           3,
	}

	producer, err := kafka.NewProducer(&producerConfig)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Не удалось создать продюсера: %s\n", err),
		)
	}

	defer func() {
		producer.Flush(int(time.Minute.Milliseconds()))
		producer.Close()
	}()

	log.Printf("Продюсер создан %v\n", producer)

	isRunning := true

	err = sendMessage(producer, daemonData)
	if err != nil {
		return err
	}

	for isRunning {
		select {
		case <-daemonData.ExistSignalChannel:
			isRunning = false
			continue
		case <-time.After(time.Minute):
			err := sendMessage(producer, daemonData)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("Завершение работы продюсера\n")
	daemonData.WaitGroup.Done()

	return nil
}

func sendMessage(producer *kafka.Producer, daemonData *DaemonSharedData) error {
	messageValue, err := json.Marshal(&SomeMessage{
		Message: "Hello at " + time.Now().String(),
	})
	if err != nil {
		return errors.New(
			fmt.Sprintf("Не удалось сериализовать сообщение: %s\n", err),
		)
	}

	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &TopicName,
			Partition: kafka.PartitionAny,
		},
		Value: messageValue,
	}, deliveryChan)
	if err != nil {
		return errors.New(fmt.Sprintf("Не удалось отправить сообщение: %s\n", err))
	}

	e := <-deliveryChan

	message := e.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		log.Printf("Сообщение не отправлено: %v\n", message.TopicPartition.Error)
	} else {
		log.Printf("Сообщение отправлено на %v\n", message.TopicPartition)

		daemonData.PushChannel <- message
	}

	close(deliveryChan)
	return nil
}
