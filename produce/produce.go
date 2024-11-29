package produce

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

// RabbitMQ 메시지 발행 함수
func PublishToRabbitMQ(queueName string, message string) error {
	// RabbitMQ 연결 URL
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/" // 기본값 설정
	}

	// RabbitMQ 연결 설정
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Printf("RabbitMQ 연결 실패: %v", err)
		return err
	}
	defer conn.Close()

	// 채널 생성
	channel, err := conn.Channel()
	if err != nil {
		log.Printf("RabbitMQ 채널 생성 실패: %v", err)
		return err
	}
	defer channel.Close()

	// 큐 선언
	_, err = channel.QueueDeclare(
		queueName, // 큐 이름
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("큐 선언 실패: %v", err)
		return err
	}

	// 메시지 발행
	err = channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("메시지 발행 실패: %v", err)
		return err
	}

	log.Printf("메시지 발행 성공: %s", message)
	return nil
}
