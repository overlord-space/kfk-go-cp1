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

	go func() {
		err := internal.CreateProducer(&daemonData)
		if err != nil {
			log.Print(err)
			close(exitSystemSignal)
		}
	}()
	time.Sleep(time.Second * 5)
	log.Printf("\n")

	go func() {
		err := internal.CreatePullCustomer(&daemonData)
		if err != nil {
			log.Print(err)
			close(exitSystemSignal)
		}
	}()
	go func() {
		err := internal.CreatePushConsumer(&daemonData)
		if err != nil {
			log.Print(err)
			close(exitSystemSignal)
		}
	}()

	<-exitSystemSignal
	log.Printf("Получен сигнал завершения работы\n")
	close(daemonData.ExistSignalChannel)
	close(daemonData.PushChannel)

	daemonData.WaitGroup.Wait()
}
