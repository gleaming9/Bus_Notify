package model

// AlertMessage 구조체: RabbitMQ에서 수신할 메시지의 구조
type AlertMessage struct {
	Email   string `json:"email"`   // 수신자 이메일
	Subject string `json:"subject"` // 이메일 제목
	Body    string `json:"body"`    // 이메일 내용
}

// AlertRequest는 사용자로부터 받을 데이터 구조입니다.
type AlertRequest struct {
	StationName string `json:"stationName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	TargetTime  string `json:"targetTime" binding:"required"` // HH:mm 형식
}
