# Go Message queue examples

Go Message Queue examples (ActiveMQ, Kafka, RabbitMQ).

## Description

App contains simple req/res service which sends and receives messages through popular message queues.

## Structure

- `main` - main file to initialize and run message queues
- `config/.env` - contains environment variables
- `deploy/docker-compose.yml` - docker compose file with application dependencies
- `internal/activemq` - contains ActiveMQ req/res implementation
- `internal/rabbitmq` - contains RabbitMQ req/res implementation
- `internal/kafka` - contains Kafka req/res implementation
- `internal/mq` - contains MQ interface
- `internal/config` - parses and provides application configuration data

## Running

1. `cd ./deploy`
2. `docker-compose up -d`
3. `cd ..`
4. `go run ./cmd/main.go`
