package services

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/streadway/amqp"
)

// 메시지 발행
func PublishToRabbitMQ(ch *amqp.Channel, routingKey, message string) {
	err := ch.Publish(
		"bus_exchange", // 교환기 이름
		routingKey,     // 라우팅 키
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("RabbitMQ 메시지 발행 실패: %v", err)
	}
	fmt.Println("알림 메시지 발행 성공:", message)
}

func SendEmail(to, subject, body string) {
	from := "yoonaji@khu.ac.kr"
	password := "wldbsdk3895"
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatalf("이메일 전송 실패: %v", err)
	}

	log.Println("이메일 전송 성공:", to)
}
