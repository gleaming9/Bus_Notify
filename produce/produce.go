package produce

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gleaming9/Bus_Notify/model"
	"github.com/streadway/amqp"
)

func PublishToRabbitMQ(queueName string, message model.AlertMessage) error {
	log.Printf("RabbitMQ 연결 준비, 큐 이름: %s", queueName)

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("메시지 직렬화 실패: %v", err)
		return err
	}

	// RabbitMQ 연결 URL
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/" // 기본값 설정
	}

	// RabbitMQ 연결
	conn, err := amqp.Dial(rabbitMQURL)
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
			ContentType: "application/json",
			Body:        jsonMessage,
		},
	)
	if err != nil {
		log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
		return err
	}

	log.Printf("RabbitMQ 메시지 발행 성공: %s", string(jsonMessage))
	return nil
}
