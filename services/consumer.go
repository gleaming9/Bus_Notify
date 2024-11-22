package services

import (
	"log"

	"github.com/streadway/amqp"
)

func ReceiveFromRabbitMQ(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // 큐 이름
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("RabbitMQ 메시지 소비 실패: %v", err)
	}

	// 메시지 처리
	for msg := range msgs {
		log.Printf("받은 메시지: %s", msg.Body)
		// 여기에 사용자에게 알림을 보내는 로직 추가
		SendEmail("user@example.com", "버스 알림", string(msg.Body))
	}
}
