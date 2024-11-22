package outputs

import (
	"fmt"
	"github.com/gleaming9/Bus_Notify/api"
	"log"
)

func PrintBusInfo(stationID string) {
	// main 함수는 테스트 목적으로 제공됩니다.
	// serviceKey에 발급받은 서비스 키(인코딩 키)를 입력하고, stationID에 정류소 ID를 입력한 후 실행해 보면 정상적으로 동작합니다
	serviceKey := "FeGUV3k8vqkcTd05EuMbi%2F0kjfLT7YbRP2xwHUxPqrZyBVJGfVA5lfvyIhqOKL1%2FYU5tbctcadarl5Jj3Ym4vg%3D%3D" // 공공데이터포털에서 발급받은 서비스 키를 입력하세요.
	//stationID := "200000078"                                                                                         // 예시 정류소 ID

	// 버스 도착 정보 가져오기
	result, err := api.GetBusArrivalInfo(serviceKey, stationID)
	if err != nil {
		log.Fatalf("API 호출 실패: %v", err)
	}

	// 버스 도착 정보 출력
	fmt.Printf("\n1번째 버스 도착 정보\n")
	for _, bus := range result.Body.BusArrivalList {
		fmt.Printf("노선 ID: %s, 도착 예상 시간: %s분, 차량 번호: %s, 빈자리: %s\n",
			bus.RouteID, bus.PredictTime1, bus.PlateNo1, bus.RemainSeatCnt1)
	}
	fmt.Printf("\n2번째 버스 도착 정보\n")
	for _, bus := range result.Body.BusArrivalList {
		fmt.Printf("노선 ID: %s, 도착 예상 시간: %s분, 차량 번호: %s, 빈자리: %s\n",
			bus.RouteID, bus.PredictTime2, bus.PlateNo2, bus.RemainSeatCnt2)
	}
}
