package model

// AlertMessage 구조체: RabbitMQ에서 수신할 메시지의 구조
type AlertMessage struct {
	StationName string `json:"stationName"` // 정류소 이름
	Email       string `json:"email"`       // 수신자 이메일
	Subject     string `json:"subject"`     // 이메일 제목
	Body        string `json:"body"`        // 이메일 내용
	TargetTime  string `json:"targetTime"`  // 도착 시간
}

// AlertRequest 는 사용자로부터 받을 데이터 구조입니다.
type AlertRequest struct {
	StationName string `json:"stationName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	TargetTime  string `json:"targetTime" binding:"required"` // HH:mm 형식
}

// BusData 구조체 정의
type BusData struct {
	First  string // 버스 이름
	Second string // 예상 도착 시간
}
