package main

import (
	"fmt"
	"github.com/gleaming9/Bus_Notify/api"
	"github.com/gleaming9/Bus_Notify/outputs"
	"log"
)

func main() {
	//정류소 데이터 초기화
	err := api.LoadStationData()
	if err != nil {
		log.Fatalf("정류소 데이터 로드 실패: %v", err)
	}

	for { // 시작 부분
		stationName := ""
		fmt.Printf("정류소명을 입력하세요: ")
		fmt.Scanf("%s", &stationName)

		stationID, err := api.GetStationID(stationName) // 정류소명을 입력받아 정류소 ID 출력
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		//fmt.Printf("정류소명 '%s'의 ID는 '%s'입니다.\n", stationName, stationID)

		outputs.PrintBusInfo(stationID) // 버스 도착 정보 출력
		fmt.Println()
	}
}
