package consume

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/streadway/amqp"
)

// 이메일 전송 함수
func sendEmail(to, subject, body string) {
	from := "yunaji0824@gmail.com"
	password := "auej ihcc naxw lsao"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 이메일 메시지 작성
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// SMTP 연결 및 이메일 전송
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Fatalf("이메일 전송 실패: %v", err)
	}

	fmt.Println("이메일 전송 성공!")
}

// RabbitMQ 소비자 함수 (내보내기 함수)
func ConsumeFromRabbitMQ() {
	// RabbitMQ 연결
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	log.Println("RabbitMQ 메시지 수신 대기 중...")

	// 메시지 처리
	for message := range messages {
		log.Printf("메시지 수신: %s", message.Body)

		// 메시지 내용을 이메일로 전송
		sendEmail("yoonaji@khu.ac.kr", "버스 알림 서비스", string(message.Body))
	}
}
