package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gleaming9/Bus_Notify/api"
	"github.com/gleaming9/Bus_Notify/consume"
	"github.com/gleaming9/Bus_Notify/mail"
	"github.com/gleaming9/Bus_Notify/outputs"
)

func main() {
	// 정류소 데이터 초기화
	err := api.LoadStationData()
	if err != nil {
		log.Fatalf("정류소 데이터 로드 실패: %v", err)
	}

	for { // 시작 부분
		stationName := ""
		fmt.Printf("정류소명을 입력하세요: ")
		fmt.Scanf("%s", &stationName)

		// GetStationID로 정류소 ID 조회
		stationID, err := api.GetStationID(stationName)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		busArrivals, err := outputs.GetBusInfo(stationID)
		if err != nil {
			log.Printf("버스 도착 정보 가져오기 실패: %v", err)
			continue
		}

		// 도착 정보를 출력
		for _, bus := range busArrivals {
			fmt.Printf("버스 번호: %s, 도착 시간: %s분, 남은 좌석: %s\n",
				bus.BusNumber, bus.ArrivalTime, bus.RemainSeats)
		}

		// RabbitMQ 메시지 발행 및 소비
		var wg sync.WaitGroup
		wg.Add(2)

		// mail.go 실행
		go func() {
			defer wg.Done()
			mail.GenerateAndPublishBusAlerts(stationName, busArrivals)
		}()

		// consume.go 실행
		go func() {
			defer wg.Done()
			consume.ConsumeFromRabbitMQ()
		}()

		wg.Wait()
	}
}
