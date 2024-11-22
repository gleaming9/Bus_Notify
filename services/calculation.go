package services

import (
	"bus_notify/models"
	"log"
	"time"
)

// 출발 시간 계산
func CalculateDepartureTime(userInput models.UserInput, arrival models.BusArrival) string {
	targetTime, err := time.Parse("15:04", userInput.TargetTime)
	if err != nil {
		log.Fatalf("목표 시간 파싱 실패: %v", err)
	}

	// 버스 도착 시간 계산
	firstBusArrival := targetTime.Add(-time.Duration(arrival.PredictTime1+userInput.WalkingTime) * time.Minute)
	secondBusArrival := targetTime.Add(-time.Duration(arrival.PredictTime2+userInput.WalkingTime) * time.Minute)

	// 가장 빠른 출발 시간 반환
	if arrival.PredictTime2 > 0 && secondBusArrival.After(firstBusArrival) {
		return firstBusArrival.Format("15:04")
	}
	return secondBusArrival.Format("15:04")
}
