package api

import (
	"encoding/csv"
	"fmt"
	"os"
)

// StationMap: 정류소명을 키로, 정류소 ID를 값으로 저장하는 맵
var StationMap map[string]string

// LoadStationData: CSV 파일에서 정류소 데이터를 로드하여 StationMap을 초기화
func LoadStationData() error {
	// CSV 파일 열기
	filePath := "bus_stations.csv"
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("CSV 파일 열기 실패: %v", err)
	}
	defer file.Close()

	// CSV 데이터 읽기
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("CSV 데이터 읽기 실패: %v", err)
	}

	// StationMap 초기화
	StationMap = make(map[string]string)

	// 데이터 파싱
	for _, record := range records {
		// 정류소명 (예: "가평역")과 정류소 ID (예: "239000818") 추출
		if len(record) < 2 {
			continue // 데이터가 부족하면 스킵
		}
		stationName := record[1] // 정류소명 (두 번째 열)
		stationID := record[3]   // 정류소 ID (네 번째 열)

		StationMap[stationName] = stationID
	}
	return nil
}

// GetStationID: 정류소명으로 정류소 ID를 검색하는 함수
func GetStationID(stationName string) (string, error) {
	if id, exists := StationMap[stationName]; exists {
		return id, nil
	}
	return "", fmt.Errorf("정류소 '%s'에 대한 ID를 찾을 수 없습니다", stationName)
}
