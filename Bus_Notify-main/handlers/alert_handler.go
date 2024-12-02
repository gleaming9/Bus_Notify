package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gleaming9/Bus_Notify/model"
	"github.com/gleaming9/Bus_Notify/service"
	"net/http"
	"time"
)

// 입력된 시간 파싱 및 KST로 설정
func parseInputTime(targetTimeStr string) (time.Time, error) {
	// 시간 문자열 파싱 (시간만)
	parsedTime, err := time.Parse("15:04", targetTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("시간 형식이 올바르지 않습니다 (HH:mm 형식): %v", err)
	}

	// 현재 서버 시간 가져오기
	now := time.Now()

	// 타임존을 KST로 설정
	kst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, fmt.Errorf("KST 타임존 설정 실패: %v", err)
	}

	// 현재 날짜와 파싱된 시간 병합
	finalTime := time.Date(
		now.Year(), now.Month(), now.Day(), // 현재 날짜
		parsedTime.Hour(), parsedTime.Minute(), 0, // 입력된 시간
		0, kst, // KST 타임존 설정
	)

	return finalTime, nil
}

// AlertHandler는 알림 요청을 처리하는 핸들러입니다.
func AlertHandler(c *gin.Context) {
	var req model.AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 본문입니다", "details": err.Error()})
		return
	}

	// 입력된 시간 파싱
	targetTime, err := parseInputTime(req.TargetTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "시간 형식이 올바르지 않습니다 (HH:mm 형식)"})
		return
	}

	// 현재 시간과 비교
	if targetTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "설정된 시간이 이미 지났습니다"})
		return
	}

	//service 호출
	if err := service.StartAlertService(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "알림 요청을 처리하는 중 오류가 발생했습니다"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "알림 요청이 성공적으로 등록되었습니다.",
		"station": req.StationName,
		"email":   req.Email,
		"time":    req.TargetTime,
	})
}
