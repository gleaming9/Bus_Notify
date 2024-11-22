package models

// 버스 도착 정보를 표현하는 구조체
type BusArrival struct {
	StationID    string `json:"stationId"`
	BusRouteID   string `json:"routeId"`
	PredictTime1 int    `json:"predictTime1"` // 첫 번째 도착 예상 시간 (분 단위)
	PredictTime2 int    `json:"predictTime2"` // 두 번째 도착 예상 시간 (분 단위)
}

// 사용자 입력 데이터를 표현하는 구조체
type UserInput struct {
	TargetTime  string // 목표 시간 (예: "08:40")
	StationID   string // 정류소 ID
	BusRouteID  string // 버스 노선 ID
	WalkingTime int    // 집에서 정류소까지 소요 시간 (분 단위)
}
