package internal

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"sync"
)

type SomeMessage struct {
	Message string
}

type ExitSignalChannel chan struct{}

type DaemonSharedData struct {
	ExistSignalChannel ExitSignalChannel
	WaitGroup          *sync.WaitGroup
	PushChannel        chan *kafka.Message
}
