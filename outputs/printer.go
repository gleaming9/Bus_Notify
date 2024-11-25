package outputs

import (
	"fmt"
	"log"

	"github.com/gleaming9/Bus_Notify/api"
)

// PrintBusInfo는 정류소 ID를 기반으로 버스 도착 정보를 출력하는 함수입니다.
func PrintBusInfo(stationID string) {
	serviceKey := "FeGUV3k8vqkcTd05EuMbi%2F0kjfLT7YbRP2xwHUxPqrZyBVJGfVA5lfvyIhqOKL1%2FYU5tbctcadarl5Jj3Ym4vg%3D%3D" // 발급받은 서비스 키

	// 버스 도착 정보 가져오기
	result, err := api.GetBusArrivalInfo(serviceKey, stationID)
	if err != nil {
		log.Fatalf("API 호출 실패: %v", err)
	}

	// 1번째 버스 도착 정보 출력
	fmt.Printf("\n1번째 버스 도착 정보\n")
	for _, bus := range result.Body.BusArrivalList {
		routeInfo, err := api.GetBusRouteInfo(serviceKey, bus.RouteID) // 노선 ID를 기반으로 버스 노선 정보 가져오기
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
		routeInfo, err := api.GetBusRouteInfo(serviceKey, bus.RouteID) // 노선 ID를 기반으로 버스 노선 정보 가져오기
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

// ///여기부터추가함
type BusArrival struct {
	BusNumber   string
	ArrivalTime string
	RemainSeats string
}

// GetBusInfo 함수는 정류소 ID를 기반으로 실제 버스 도착 정보를 반환합니다.
func GetBusInfo(stationID string) ([]BusArrival, error) {
	serviceKey := "FeGUV3k8vqkcTd05EuMbi%2F0kjfLT7YbRP2xwHUxPqrZyBVJGfVA5lfvyIhqOKL1%2FYU5tbctcadarl5Jj3Ym4vg%3D%3D" // 발급받은 서비스 키

	// 버스 도착 정보 가져오기
	result, err := api.GetBusArrivalInfo(serviceKey, stationID)
	if err != nil {
		log.Printf("API 호출 실패: %v", err)
		return nil, err
	}

	var arrivals []BusArrival

	// 첫 번째 버스 도착 정보를 배열에 추가
	for _, bus := range result.Body.BusArrivalList {
		routeInfo, err := api.GetBusRouteInfo(serviceKey, bus.RouteID) // 노선 ID를 기반으로 버스 노선 정보 가져오기
		if err != nil {
			log.Printf("노선 ID %s의 추가 정보를 가져오는 데 실패했습니다: %v\n", bus.RouteID, err)
			continue
		}

		arrivals = append(arrivals, BusArrival{
			BusNumber:   routeInfo.MsgBody.BusRouteInfoItem.RouteName,
			ArrivalTime: bus.PredictTime1,
			RemainSeats: bus.RemainSeatCnt1,
		})
	}

	// 두 번째 버스 도착 정보를 배열에 추가
	for _, bus := range result.Body.BusArrivalList {
		routeInfo, err := api.GetBusRouteInfo(serviceKey, bus.RouteID)
		if err != nil {
			log.Printf("노선 ID %s의 추가 정보를 가져오는 데 실패했습니다: %v\n", bus.RouteID, err)
			continue
		}

		arrivals = append(arrivals, BusArrival{
			BusNumber:   routeInfo.MsgBody.BusRouteInfoItem.RouteName,
			ArrivalTime: bus.PredictTime2,
			RemainSeats: bus.RemainSeatCnt2,
		})
	}

	return arrivals, nil
}
