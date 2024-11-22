package main

import (
	"fmt"
	"log"
)

func main() {
	// CSV 파일 경로
	csvFilePath := "bus_stations.csv"

	// CSV 데이터 로드
	err := api.LoadStationData(csvFilePath)
	if err != nil {
		log.Fatalf("정류소 데이터 로드 실패: %v", err)
	}

	// 정류소명을 입력받아 정류소 ID 출력
	stationName := "가평역" // 예: 사용자가 검색하는 정류소명
	stationID, err := api.GetStationID(stationName)
	if err != nil {
		log.Fatalf("정류소 ID 검색 실패: %v", err)
	}
	fmt.Printf("정류소명 '%s'의 ID는 '%s'입니다.\n", stationName, stationID)
}
