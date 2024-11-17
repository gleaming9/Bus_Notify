package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const baseURL = "http://apis.data.go.kr/6410000/busarrivalservice/getBusArrivalList"

// API 응답 구조체 정의
type BusArrivalListResponse struct {
	XMLName xml.Name `xml:"response"`
	Header  struct {
		QueryTime     string `xml:"queryTime"`
		ResultCode    int    `xml:"resultCode"`
		ResultMessage string `xml:"resultMessage"`
	} `xml:"msgHeader"`
	Body struct {
		BusArrivalList []struct {
			Flag           string `xml:"flag"`
			LocationNo1    string `xml:"locationNo1"`
			LocationNo2    string `xml:"locationNo2"`
			LowPlate1      string `xml:"lowPlate1"`
			LowPlate2      string `xml:"lowPlate2"`
			PlateNo1       string `xml:"plateNo1"`
			PlateNo2       string `xml:"plateNo2"`
			PredictTime1   string `xml:"predictTime1"`
			PredictTime2   string `xml:"predictTime2"`
			RemainSeatCnt1 string `xml:"remainSeatCnt1"`
			RemainSeatCnt2 string `xml:"remainSeatCnt2"`
			RouteID        string `xml:"routeId"`
			StaOrder       string `xml:"staOrder"`
			StationID      string `xml:"stationId"`
		} `xml:"busArrivalList"`
	} `xml:"msgBody"`
}

// GetBusArrivalInfo: 정류소 ID를 기반으로 버스 도착 정보를 가져오는 함수
func GetBusArrivalInfo(serviceKey, stationID string) (*BusArrivalListResponse, error) {
	// URL 구성
	params := url.Values{}
	params.Add("stationId", stationID)
	fullURL := fmt.Sprintf("%s?serviceKey=%s&%s", baseURL, serviceKey, params.Encode())

	log.Println("Request URL:", fullURL)

	// HTTP GET 요청
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 응답 실패: HTTP %d", resp.StatusCode)
	}

	// 응답 본문 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 본문 읽기 실패: %v", err)
	}

	// XML 데이터 파싱
	var response BusArrivalListResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("XML 파싱 실패: %v", err)
	}

	// 결과 코드 확인
	if response.Header.ResultCode != 0 {
		return nil, fmt.Errorf("API 오류: %s", response.Header.ResultMessage)
	}

	return &response, nil
}

//메인 실행 파일

// main 함수는 테스트 목적으로 제공됩니다.
// serviceKey에 발급받은 서비스 키(인코딩 키)를 입력하고, stationID에 정류소 ID를 입력한 후 실행해 보면 정상적으로 동작합니다
func main() {
	serviceKey := "FeGUV3k8vqkcTd05EuMbi%2F0kjfLT7YbRP2xwHUxPqrZyBVJGfVA5lfvyIhqOKL1%2FYU5tbctcadarl5Jj3Ym4vg%3D%3D" // 공공데이터포털에서 발급받은 서비스 키를 입력하세요.
	stationID := "200000078"                                                                                         // 예시 정류소 ID

	// 버스 도착 정보 가져오기
	result, err := GetBusArrivalInfo(serviceKey, stationID)
	if err != nil {
		log.Fatalf("API 호출 실패: %v", err)
	}

	// 결과 출력
	for _, bus := range result.Body.BusArrivalList {
		fmt.Printf("노선 ID: %s, 도착 예상 시간: %s분, 차량 번호: %s, 빈자리: %s\n",
			bus.RouteID, bus.PredictTime1, bus.PlateNo1, bus.RemainSeatCnt1)
	}
}
