package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

// GetBusArrivalInfo : 정류소 ID를 기반으로 버스 도착 정보를 가져오는 함수
func GetBusArrivalInfo(serviceKey, stationID string) (*BusArrivalListResponse, error) {
	// URL 구성
	params := url.Values{}
	params.Add("stationId", stationID)
	fullURL := fmt.Sprintf("%s?serviceKey=%s&%s", baseURL, serviceKey, params.Encode())

	//log.Println("Request URL:", fullURL)

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
