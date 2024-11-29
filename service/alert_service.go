package service

import (
	"fmt"
	"github.com/gleaming9/Bus_Notify/api"
	"github.com/gleaming9/Bus_Notify/model"
	"github.com/gleaming9/Bus_Notify/produce"
	"log"
	"os"
	"strconv"
	"time"
)

// StartAlertService: 정류소 이름과 목표 시간을 받아 알림 모니터링을 시작
func StartAlertService(request *model.AlertRequest) error {
	stationName := request.StationName

	go func() {
		alertSent := [3]bool{false, false, false} // 15분, 10분, 5분
		for {
			// 정류소 ID 가져오기
			stationID, err := api.GetStationID(stationName)
			if err != nil {
				log.Printf("정류소 ID 조회 실패: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			// 버스 도착 정보 가져오기
			arrivalInfo, err := api.GetBusArrivalInfo(os.Getenv("SERVICE_KEY"), stationID)
			if err != nil {
				log.Printf("버스 도착 정보 갱신 실패: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			// 알림 조건 체크 및 RabbitMQ 메시지 발행
			checkAndSendAlerts(arrivalInfo, stationName, &alertSent)

			// 모든 알림이 전송되면 종료
			if alertSent[0] && alertSent[1] && alertSent[2] {
				log.Println("모든 알림이 전송 완료되었습니다.")
				break
			}

			// 1분 대기
			time.Sleep(1 * time.Minute)
		}
	}()
	return nil
}

func checkAndSendAlerts(arrivalInfo *api.BusArrivalListResponse, stationName string, alertSent *[3]bool) {
	// 알림 조건 (15분, 10분, 5분 남았을 때 알림 전송)
	thresholds := []int{15, 10, 5}

	firstArrival := arrivalInfo.Body.BusArrivalList[0]

	// PredictTime1을 정수로 변환
	predictTime, err := strconv.Atoi(firstArrival.PredictTime1)
	if err != nil {
		log.Printf("PredictTime1 변환 실패: %v", err)
		return
	}

	for i, threshold := range thresholds {
		// 해당 조건의 알림이 아직 전송되지 않았고, 남은 시간이 threshold 이하일 경우
		if !alertSent[i] && predictTime <= threshold {
			// 메시지 작성
			message := buildAlertMessage(stationName, firstArrival.PlateNo1, threshold)

			// RabbitMQ로 메시지 전송
			err := produce.PublishToRabbitMQ("bus_alerts", message)
			if err != nil {
				log.Printf("RabbitMQ 메시지 전송 실패: %v", err)
				continue
			}

			// 알림 상태 업데이트
			alertSent[i] = true
			log.Printf("알림 전송 완료: %s", message)
		}
	}
}

func buildAlertMessage(stationName, busNumber string, timeLeft int) string {
	return fmt.Sprintf("%s 정류소에서 %s 버스가 %d분 후 도착합니다.", stationName, busNumber, timeLeft)
}
