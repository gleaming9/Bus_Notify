package main

import (
	"bus_notify/config"
	"bus_notify/models"
	"bus_notify/services"
	"fmt"
	"log"
)

func main() {
	// RabbitMQ 연결 및 채널 설정
	conn, ch := config.ConnectToRabbitMQ()
	defer conn.Close()
	defer ch.Close()

	// 교환기 선언
	err := config.DeclareExchange(ch, "bus_exchange")
	if err != nil {
		log.Fatalf("RabbitMQ 교환기 선언 실패: %v", err)
	}

	// 사용자 입력 데이터
	userInput := models.UserInput{
		TargetTime:  "08:40",
		StationID:   "200000078",
		BusRouteID:  "200000085",
		WalkingTime: 10, // 10분
	}

	// API 호출
	arrival, err := services.FetchBusArrivalData(userInput.StationID, userInput.BusRouteID)
	if err != nil {
		log.Fatalf("버스 API 호출 실패: %v", err)
	}

	// 출발 시간 계산
	departureTime := services.CalculateDepartureTime(userInput, *arrival)
	message := fmt.Sprintf("집에서 %s에 출발해야 합니다.", departureTime)

	// 메시지 발행
	go services.PublishToRabbitMQ(ch, "bus.notification", message)
	select {}
}
