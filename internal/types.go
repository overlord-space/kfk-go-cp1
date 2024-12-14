package internal

import (
	"sync"
)

type SomeMessage struct {
	Message string
}

type ExitSignalChannel chan struct{}

type DaemonSharedData struct {
	ExistSignalChannel ExitSignalChannel
	WaitGroup          *sync.WaitGroup
}
