package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	stationAPIBaseURL = "http://apis.data.go.kr/6410000/busarrivalservice/getBusArrivalList"
	serviceKey        = "8IK6yQqxcJbxtB0xPTIJIqIK5Cc%2FyRW9AolNym1VpDyP3k8egXaargvK57AjmSO3FaLnJ0m%2F7hgidr38niwyZA%3D%3D"
)

type Station struct {
	StationID   string `json:"stationId"`
	StationName string `json:"stationName"`
}

type SOAPFault struct {
	Fault Fault `xml:"Body>Fault"`
}

type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}

func FetchStationID(stationName string) ([]Station, error) {
	req, err := http.NewRequest("GET", stationAPIBaseURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("serviceKey", serviceKey)
	query.Add("keyword", stationName) // 검색 키워드: 정류소 이름
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("API 응답 내용: %s", string(body)) // 응답 내용 출력

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 호출 실패, 상태 코드: %d, 응답 내용: %s", resp.StatusCode, string(body))
	}

	// JSON 응답 파싱
	var result struct {
		Response struct {
			Body struct {
				Items struct {
					Item []Station `json:"item"`
				} `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		// XML 응답 처리
		var soapFault SOAPFault
		if xmlErr := xml.Unmarshal(body, &soapFault); xmlErr == nil {
			return nil, fmt.Errorf("API 오류: %s - %s", soapFault.Fault.FaultCode, soapFault.Fault.FaultString)
		}

		return nil, fmt.Errorf("JSON 파싱 실패: %v\n응답 내용: %s", err, string(body))
	}

	return result.Response.Body.Items.Item, nil
}

func main() {
	stationName := "경희대정문" // 검색하고 싶은 정류소 이름
	stations, err := FetchStationID(stationName)
	if err != nil {
		log.Fatalf("정류소 검색 실패: %v", err)
	}

	for _, station := range stations {
		fmt.Printf("정류소 ID: %s, 정류소 이름: %s\n", station.StationID, station.StationName)
	}
}
