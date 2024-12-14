package internal

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"time"
)

func CreateProducer(daemonData *DaemonSharedData) *kafka.Producer {
	producerConfig := kafka.ConfigMap{
		"bootstrap.servers": BootstrapServers,
	}

	producer, err := kafka.NewProducer(&producerConfig)
	if err != nil {
		log.Fatalf("Не удалось создать продюсера: %s\n", err)
	}
	defer producer.Close()

	log.Printf("Продюсер создан %v\n", producer)

	isRunning := true

	sendMessage(producer)

	for isRunning {
		select {
		case <-daemonData.ExistSignalChannel:
			isRunning = false
			continue
		case <-time.After(time.Minute):
			sendMessage(producer)
		}
	}

	log.Printf("Завершение работы продюсера\n")
	daemonData.WaitGroup.Done()

	return producer
}

func sendMessage(producer *kafka.Producer) {
	messageValue, err := json.Marshal(&SomeMessage{
		Message: "Hello at " + time.Now().String(),
	})
	if err != nil {
		log.Fatalf("Не удалось сериализовать сообщение: %s\n", err)
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
		log.Fatalf("Не удалось отправить сообщение: %s\n", err)
	}

	e := <-deliveryChan

	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.Printf("Сообщение не отправлено: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Сообщение отправлено на %v\n", m.TopicPartition)
	}

	close(deliveryChan)
}
