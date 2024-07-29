#  Микросервис обработки сообщений
Тестовое задание на вакансию "Junior Go-разработчик"

### Задание:
Разработать микросервис на Go, который будет:
* принимать сообщения через HTTP API,
* сохранять их в PostgreSQL,
* а затем отправлять в Kafka для дальнейшей обработки.
* Обработанные сообщения должны помечаться.
* Сервис должен также предоставлять API для получения статистики по обработанным сообщениям.

### Требования:
	1.	Использовать Go 1.20+
	2.	Использовать PostgreSQL для хранения сообщений
	3.	Реализовать отправку и чтение сообщений в Kafka
	4.	Предоставить возможность запустить проект в Docker

## Решение

Приложение реализовано в соответствии с приведенными требованиями. Для запуска приложения необходимы:
* Docker Engine версии >= 24
* Docker Compose >= v2.24.0
* Доступ к интернету для скачивания docker images

### Доступ к приложению в интернет
Приложение принимает сообщения в следующем формате JSON по **HTTP API через uri /api**
```
{
   "msg": "Текст сообщения"
}
```
* **Отправка сообщения в приложение**
```
curl -i --request POST \
  --header 'Authorization: sjdflks9r4.-Sdjf' \
  --header 'Content-Type: application/json' \
  --data '{"msg":"test message"}' \
  --url http://87.242.87.183:8080/api
```

С целью защиты приложения от несанкционированных подключений я добавил токен аутентификации, который передается в заголовке **Authorization** , в приведенном примере данный заголовок указан с **работающим токеном**.

Просмотреть метрики статистики полученных и обработанных сообщений можно через **HTTP API uri /metrics**
* **Статистика обработанных сообщений**
```
curl -i --request GET \
  --header 'Authorization: sjdflks9r4.-Sdjf' \
  --url http://87.242.87.183:8080/metrics 
```
ответ будет в виде JSON сообщения вида (пример):
```
{
	"total": 13,
	"processed": 13
}
```

### Запуск локально в Docker

В репозитории **уже** лежит конфигурационный файл **infra/.env** он необходим для запуска проекта в Docker.

Для запуска приложения в docker как есть необходимо:
* в консоли linux перейти в каталог infra, в каталоге приложения:
```
cd msg-proc/infra
```
* запустить проект в docker
```
docker-compose up -d
```
в результате скачаются необходимые docker images:
* confluentinc/cp-zookeeper:7.7.0
* confluentinc/cp-kafka:7.7.0
* provectuslabs/kafka-ui:v0.7.2
* postgres:16.3-alpine
* golang:1.22.5
* alpine:3.20

запуститься необходимая инфраструктура включающая в себя СУБД postgreSQL, kafka, kafka-ui, zookeeper

и соберется приложение в локальный docker image **msg-proc** при этом промежуточные слои можно удалить запустив 
```
docker system prune -f
```
После этого приложение будет доступно по адресу сервера/компьютера на котором выполняли команду запуска, а так же локально через loopback на порту **8080**.


Файл **.env** содержит следующие переменные:
```
POSTGRES_PASSWORD="postgres"
ADDRESS=0.0.0.0:8080
DATABASE_DSN="postgres://postgres:${POSTGRES_PASSWORD}@db:5432/postgres"
KAFKA_BROKERS="kafka:29092"
KAFKA_TOPIC="messages"
```
Где:
* POSTGRES_PASSWORD - пароль пользователя postgres необходимый для запуска контейнера
БД postgresql
* ADDRESS - адрес на котором приложение будет прослушивать подключения в формате ip:port
в данном случае необходимо использовать 0.0.0.0:8080
* DATABASE_DSN - строка подключения приложения к БД postgreSQL в формате "postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/DB_NAME"
* KAFKA_BROKERS - адреса kafka хостов через запятую (если их несколько в кластере) в формате "BROKER_HOST:BROKER_PORT,BROKER_HOST:BROKER_PORT,..."
* KAFKA_TOPIC - имя kafka топика в который приложение будет посылать сообщения.

В текущем решении подключение к kafka без аутентификации и без использования защиты подключения. Таблица в базе создается автоматически, топик в kafka так же создается автоматически.


### Подключение к приложению локально

Приложение принимает сообщения в следующем формате JSON по **HTTP API через uri /api**
```
{
   "msg": "Текст сообщения"
}
```
```
curl -i --request POST \
  --header 'Authorization: sjdflks9r4.-Sdjf' \
  --header 'Content-Type: application/json' \
  --data '{"msg":"test message"}' \
  --url http://localhost:8080/api 
```

Просмотреть метрики статистики полученных и обработанных сообщений можно через **HTTP API uri /metrics**
```
curl -i --request GET \
  --header 'Authorization: sjdflks9r4.-Sdjf' \
  --url http://localhost:8080/metrics
```
ответ будет в виде JSON сообщения вида:
```
{
	"total": 13,
	"processed": 13
}
```
Так же можно подключиться из браузера по адресу **localhost:8090** и посмотреть графический интерфейс kafka, реализованный приложением kafka-ui

## Собрать приложение вне контейнера
### Запуск тестов
```
go test ./...
```
### Сборка
```
go build -ldflags "-X 'main.BuildVersion=$(cat VERSION)'\
   -X 'main.BuildDate=$(date +'%Y-%m-%d %H:%M')'\
   -X 'main.BuildCommit=$(git log --oneline -1|awk '{print $1}')'"\
   -tags netgo,osusergo\
   -o app cmd/app/main.go
```