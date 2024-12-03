package consume

import (
	"fmt"
	"github.com/gleaming9/Bus_Notify/model"
	"github.com/goccy/go-json"
	"log"
	"net/smtp"
	"os"

	"github.com/streadway/amqp"
)

// 이메일 전송 함수
func sendEmail(to, subject, body string) error {
	from := "mathasdf0@gmail.com"
	password := "vxrq llox tohy brca"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 이메일 메시지 작성
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// SMTP 연결 및 이메일 전송
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Fatalf("이메일 전송 실패: %v", err)
		return err
	}

	fmt.Println("이메일 전송 성공!")
	return nil
}

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
		var alert model.AlertMessage
		if err := json.Unmarshal(message.Body, &alert); err != nil {
			log.Printf("메시지 디코딩 실패: %v", err)
			continue
		}

		// 이메일 전송
		if err := sendEmail(alert.Email, alert.Subject, alert.Body); err != nil {
			log.Printf("이메일 전송 실패: %v", err)
			continue
		}
	}
}
