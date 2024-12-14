package main

import (
	"checkpoint-project-m1/internal"
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
	}

	exitSystemSignal := make(chan os.Signal, 1)
	signal.Notify(exitSystemSignal, os.Interrupt)

	daemonData.WaitGroup.Add(2)

	go internal.CreateProducer(&daemonData)
	time.Sleep(time.Second * 5)
	go internal.CreatePullCustomer(&daemonData)

	<-exitSystemSignal
	log.Printf("Получен сигнал завершения работы\n")
	close(daemonData.ExistSignalChannel)

	daemonData.WaitGroup.Wait()
}
