package mail

import (
	"fmt"
	"log"

	"github.com/gleaming9/Bus_Notify/outputs"
	"github.com/streadway/amqp"
)

func GenerateAndPublishBusAlerts(stationName string, busArrivals []outputs.BusArrival) {
	// RabbitMQ 연결
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("RabbitMQ 연결 실패: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("채널 생성 실패: %v", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"bus_alerts",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("큐 선언 실패: %v", err)
	}

	//	maxCount := 3
	//	if len(busArrivals) < maxCount {
	//		maxCount = len(busArrivals) // 데이터가 3개 미만이면 전체 발행
	//	}

	// 실제 데이터를 기반으로 메시지 생성 및 발행
	//	for i := 0; i < maxCount; i++ {

	for _, bus := range busArrivals {
		message := fmt.Sprintf(
			"정류소명: %s\n버스 번호: %s\n도착 시간: %s",
			stationName, bus.BusNumber, bus.ArrivalTime,
		)

		err := channel.Publish(
			"",
			queue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		if err != nil {
			log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
		} else {
			log.Println("RabbitMQ 메시지 발행 성공:", message)
		}
	}
}

//}
