package internal

const (
	BootstrapServers = "localhost:9094"
)

// TopicName Пришлось вынести из констант, так как newProducer требует передачи указателя на переменную
var TopicName = "m1-project"
