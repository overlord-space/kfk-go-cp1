package main

import (
	"checkpoint-project-m1/internal"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	daemonData := internal.DaemonSharedData{
		ExistSignalChannel: make(internal.ExitSignalChannel),
		WaitGroup:          &sync.WaitGroup{},
		PushChannel:        make(chan *kafka.Message, 1),
	}

	exitSystemSignal := make(chan os.Signal, 1)
	signal.Notify(exitSystemSignal, os.Interrupt)

	daemonData.WaitGroup.Add(3)

	go internal.CreateProducer(&daemonData)
	time.Sleep(time.Second * 5)
	log.Printf("\n")

	go internal.CreatePullCustomer(&daemonData)
	go internal.CreatePushConsumer(&daemonData)

	<-exitSystemSignal
	log.Printf("Получен сигнал завершения работы\n")
	close(daemonData.ExistSignalChannel)
	close(daemonData.PushChannel)

	daemonData.WaitGroup.Wait()
}
