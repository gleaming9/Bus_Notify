package model

// AlertRequest는 사용자로부터 받을 데이터 구조입니다.
type AlertRequest struct {
	StationName string `json:"stationName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	TargetTime  string `json:"targetTime" binding:"required"` // HH:mm 형식
}
