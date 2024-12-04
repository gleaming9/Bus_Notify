package consume

import (
	"github.com/gleaming9/Bus_Notify/model"
	"github.com/gleaming9/Bus_Notify/service"
	"github.com/goccy/go-json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// RabbitMQ 소비자 함수 (내보내기 함수)
func ConsumeFromRabbitMQ() {
	// RabbitMQ 연결 URL
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/" // 기본값 설정
	}

	// RabbitMQ 연결
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("RabbitMQ 연결 실패: %v", err)
	}
	defer conn.Close()

	// 채널 생성
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("채널 생성 실패: %v", err)
	}
	defer channel.Close()

	// 큐 선언
	queue, err := channel.QueueDeclare(
		"bus_alerts", // 큐 이름
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("큐 선언 실패: %v", err)
	}

	// 메시지 수신
	messages, err := channel.Consume(
		queue.Name, // 큐 이름
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("메시지 소비 실패: %v", err)
	}
	// 메시지 처리
	for message := range messages {
		log.Printf("메시지 수신: %s", message.Body)

		// 메시지를 AlertMessage 구조체로 디코딩
		var req model.AlertRequest
		if err := json.Unmarshal(message.Body, &req); err != nil {
			log.Printf("메시지 디코딩 실패: %v", err)
			continue
		}

		go service.MonitorBusArrival(req)
	}
}
