Перед запуском приложения, сначала необходимо запустить Kafka.
Это можно сделать с помощью команды для `docker compose`:
```shell
docker compose -f docker-compose.kraft.yml up -d
```

Для запуска самого приложения достаточно запустить файл main.go без параметров.
Вся конфигурация приложения вынесена в файл [.env](./.env)
```shell
go run main.go
```

При запуске проекта, сперва запускается продюсер и отправляет первое сообщение в топик.
Далее он будет отправлять новое сообщение каждую минуту.
Через 5 секунд после запуска продюсера запускаются консьюмеры.
Логика каждого объекта описана в отдельных файлах:
- [Producer](internal/producer.go)
- [Pull Consumer](internal/consumer-pull.go)
- [Push Consumer](internal/consumer-push.go)

Pull консьюмер работает с помощью вызова встроенного метода `Pull` и обрабатывает сообщения в цикле.
Push консьюмер принимает сообщения через канал (от продюсера) и обрабатывает их в цикле.

Для всех объектов реализован graceful shutdown, который вызывается по сигналу `SIGTERM`.
При получении сигнала, каждый объект обрабатывает его по своему, что бы завершить работу корректно.

Сообщение в топики представляет собой структуру SomeMessage, внутри которого есть всего одно свойство `Message string`.
В него записывается текст в формате "Hello, World! %current_time%".
