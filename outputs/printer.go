package outputs

import (
	"fmt"
	"github.com/gleaming9/Bus_Notify/api"
	"log"
)

// BusArrival 구조체: 버스 도착 정보를 표현
type BusArrival struct {
	BusNumber         string `json:"busNumber"`
	FirstArrivalTime  string `json:"firstarrivalTime"`
	SecondArrivalTime string `json:"secondarrivalTime"`
}

// GetBusInfo 함수는 정류소 ID를 기반으로 첫 번째 버스 도착 정보를 반환합니다.
func GetBusInfo(stationID string) ([]BusArrival, error) {
	// 정류소 ID를 사용하여 버스 도착 정보를 가져옵니다.
	arrivalResult, err := api.GetBusArrivalInfo(stationID)
	if err != nil {
		log.Printf("API 호출 실패: %v", err)
		return nil, err
	}

	var busArrivals []BusArrival

	// 첫 번째 버스 도착 정보만 처리
	for _, bus := range arrivalResult.Body.BusArrivalList {
		// 버스 노선 정보 가져오기
		routeInfo, err := api.GetBusRouteInfo(bus.RouteID)
		if err != nil {
			log.Printf("노선 정보 조회 실패 (RouteID: %s): %v", bus.RouteID, err)
			continue
		}

		// 첫 번째 도착 정보를 추가
		busArrivals = append(busArrivals, BusArrival{
			BusNumber:         routeInfo.MsgBody.BusRouteInfoItem.RouteName,
			FirstArrivalTime:  bus.PredictTime1,
			SecondArrivalTime: bus.PredictTime2,
		})
	}

	return busArrivals, nil
}

// PrintBusInfo는 정류소 ID를 기반으로 버스 도착 정보를 출력하는 함수입니다.
func PrintBusInfo(stationID string) {
	// 버스 도착 정보 가져오기
	result, err := api.GetBusArrivalInfo(stationID)
	if err != nil {
		log.Fatalf("API 호출 실패: %v", err)
	}

	// 1번째 버스 도착 정보 출력
	fmt.Printf("\n1번째 버스 도착 정보\n")
	for _, bus := range result.Body.BusArrivalList {
		routeInfo, err := api.GetBusRouteInfo(bus.RouteID) // 노선 ID를 기반으로 버스 노선 정보 가져오기
		if err != nil {
			log.Printf("노선 ID %s의 추가 정보를 가져오는 데 실패했습니다: %v\n", bus.RouteID, err)
			continue
		}
		fmt.Printf("노선 이름: %s, 도착 예상 시간: %s분, 차량 번호: %s, 빈자리: %s\n",
			routeInfo.MsgBody.BusRouteInfoItem.RouteName,
			bus.PredictTime1,
			bus.PlateNo1,
			bus.RemainSeatCnt1,
		)
	}

	// 2번째 버스 도착 정보 출력
	fmt.Printf("\n2번째 버스 도착 정보\n")
	for _, bus := range result.Body.BusArrivalList {
		routeInfo, err := api.GetBusRouteInfo(bus.RouteID) // 노선 ID를 기반으로 버스 노선 정보 가져오기
		if err != nil {
			log.Printf("노선 ID %s의 추가 정보를 가져오는 데 실패했습니다: %v\n", bus.RouteID, err)
			continue
		}

		fmt.Printf("노선 이름: %s, 도착 예상 시간: %s분, 차량 번호: %s, 빈자리: %s\n",
			routeInfo.MsgBody.BusRouteInfoItem.RouteName,
			bus.PredictTime2,
			bus.PlateNo2,
			bus.RemainSeatCnt2,
		)
	}
}
