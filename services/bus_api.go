package services

import (
	"bus_notify/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	busAPIBaseURL = "http://apis.data.go.kr/6410000/busarrivalservice/getBusArrivalList"
	busAPIKey     = "8IK6yQqxcJbxtB0xPTIJIqIK5Cc%2FyRW9AolNym1VpDyP3k8egXaargvK57AjmSO3FaLnJ0m%2F7hgidr38niwyZA%3D%3D"
)

var httpClient = &http.Client{}

func FetchBusArrivalData(stationID, routeID string) (*models.BusArrival, error) {
	req, err := http.NewRequest("GET", busAPIBaseURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("serviceKey", busAPIKey)
	query.Add("stationId", stationID)
	query.Add("routeId", routeID)
	req.URL.RawQuery = query.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("API 응답 내용: %s", string(body))
		return nil, fmt.Errorf("API 응답 실패, 상태 코드: %d, 응답 내용: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response struct {
			Body struct {
				Items struct {
					Item []models.BusArrival `json:"item"`
				} `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Response.Body.Items.Item) == 0 {
		return nil, fmt.Errorf("버스 도착 정보가 없습니다.")
	}

	return &result.Response.Body.Items.Item[0], nil
}
