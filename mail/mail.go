package mail

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

// 더미 데이터 구조체
type BusArrival struct {
	StationName  string
	BusNumber    string
	ArrivalTime  string
	AlertMessage string
}

// RabbitMQ로 메시지 발행하는 함수
func publishToRabbitMQ(message string) error {
	// RabbitMQ 연결
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("RabbitMQ 연결 실패: %v", err)
	}
	defer conn.Close()

	// 채널 생성
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("채널 생성 실패: %v", err)
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
		return fmt.Errorf("큐 선언 실패: %v", err)
	}

	// 메시지 발행
	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("메시지 발행 실패: %v", err)
	}

	log.Println("RabbitMQ로 메시지 발행 성공:", message)
	return nil
}

// 더미 데이터 생성 함수
func generateDummyData(stationName string, busNumbers []string) []BusArrival {
	rand.Seed(time.Now().UnixNano())
	var arrivals []BusArrival

	for i := 0; i < 5; i++ {
		arrivalTime := time.Now().Add(time.Duration(rand.Intn(60)) * time.Minute)
		bus := busNumbers[rand.Intn(len(busNumbers))]
		alertMessage := fmt.Sprintf("정류소에서 %s번 버스가 %s에 도착 예정입니다.", bus, arrivalTime.Format("15:04"))
		arrivals = append(arrivals, BusArrival{
			StationName:  stationName,
			BusNumber:    bus,
			ArrivalTime:  arrivalTime.Format("15:04"),
			AlertMessage: alertMessage,
		})
	}

	return arrivals
}

// 외부에서 호출할 수 있는 함수
func GenerateAndPublishBusAlerts(stationName string, busNumbers []string) {
	// 더미 데이터 생성
	dummyData := generateDummyData(stationName, busNumbers)

	// 각 데이터를 RabbitMQ로 발행
	for _, data := range dummyData {
		message := fmt.Sprintf(
			"정류소명: %s\n버스 번호: %s\n도착 시간: %s\n메시지: %s",
			data.StationName, data.BusNumber, data.ArrivalTime, data.AlertMessage,
		)

		// RabbitMQ로 발행
		if err := publishToRabbitMQ(message); err != nil {
			log.Fatalf("RabbitMQ 발행 실패: %v", err)
		}
	}
}
