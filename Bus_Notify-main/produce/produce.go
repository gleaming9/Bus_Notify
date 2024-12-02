package produce

import (
	"log"

	"github.com/streadway/amqp"
)

func PublishToRabbitMQ(queueName string, message string) error {
	log.Printf("RabbitMQ 연결 준비, 큐 이름: %s", queueName)

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("RabbitMQ 연결 실패: %v", err)
		return err
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("RabbitMQ 채널 생성 실패: %v", err)
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("RabbitMQ 큐 선언 실패: %v", err)
		return err
	}

	log.Printf("RabbitMQ 메시지 발행: %s", message)

	err = channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
		return err
	}

	log.Printf("RabbitMQ 메시지 발행 성공: %s", message)
	return nil
}
