package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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
func GetBusArrivalInfo(stationID string) (*BusArrivalListResponse, error) {
	const baseURL = "http://apis.data.go.kr/6410000/busarrivalservice/getBusArrivalList"
	fullURL := fmt.Sprintf("%s?serviceKey=%s&stationId=%s", baseURL, os.Getenv("SERVICE_KEY"), stationID)

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
	body, err := io.ReadAll(resp.Body)
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

type BusRouteInfo struct {
	XMLName   xml.Name `xml:"response"`
	MsgHeader struct {
		QueryTime     string `xml:"queryTime"`
		ResultCode    string `xml:"resultCode"`
		ResultMessage string `xml:"resultMessage"`
	} `xml:"msgHeader"`
	MsgBody struct {
		BusRouteInfoItem struct {
			CompanyID        string `xml:"companyId"`
			CompanyName      string `xml:"companyName"`
			CompanyTel       string `xml:"companyTel"`
			DistrictCd       string `xml:"districtCd"`
			DownFirstTime    string `xml:"downFirstTime"`
			DownLastTime     string `xml:"downLastTime"`
			EndStationID     string `xml:"endStationId"`
			EndStationName   string `xml:"endStationName"`
			PeekAlloc        string `xml:"peekAlloc"`
			RegionName       string `xml:"regionName"`
			RouteID          string `xml:"routeId"`
			RouteName        string `xml:"routeName"`
			RouteTypeCd      string `xml:"routeTypeCd"`
			RouteTypeName    string `xml:"routeTypeName"`
			StartMobileNo    string `xml:"startMobileNo"`
			StartStationID   string `xml:"startStationId"`
			StartStationName string `xml:"startStationName"`
			UpFirstTime      string `xml:"upFirstTime"`
			UpLastTime       string `xml:"upLastTime"`
			NPeekAlloc       string `xml:"nPeekAlloc"`
		} `xml:"busRouteInfoItem"`
	} `xml:"msgBody"`
}

// GetBusRouteInfo 는 노선 ID를 사용하여 노선 정보를 가져오는 함수입니다.
func GetBusRouteInfo(routeID string) (*BusRouteInfo, error) {
	baseURL := "http://apis.data.go.kr/6410000/busrouteservice/getBusRouteInfoItem"
	reqURL := fmt.Sprintf("%s?serviceKey=%s&routeId=%s", baseURL, os.Getenv("SERVICE_KEY"), routeID)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 응답 오류: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", err)
	}

	var info BusRouteInfo
	if err := xml.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("XML 파싱 실패: %w", err)
	}

	return &info, nil
}
