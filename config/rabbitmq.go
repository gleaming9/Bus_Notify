package config

import (
	"log"

	"github.com/streadway/amqp"
)

const rabbitMQURL = "amqp://guest:guest@localhost:5672/"

// RabbitMQ 연결 및 채널 생성
func ConnectToRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("RabbitMQ 연결 실패: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ 채널 생성 실패: %v", err)
	}
	return conn, ch
}

// 교환기 선언
func DeclareExchange(ch *amqp.Channel, exchangeName string) error {
	return ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // args
	)
}
