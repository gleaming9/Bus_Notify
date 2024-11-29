package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gleaming9/Bus_Notify/api"
	"github.com/gleaming9/Bus_Notify/outputs"
	"log"
	"net/http"
)

// GetBusInfoHandler handles requests to get bus arrival info
func GetBusInfoHandler(c *gin.Context) {
	stationName := c.Query("stationName") // 정류소 이름 받기
	if stationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "정류소 이름이 필요합니다"})
		return
	}

	// 정류소 ID 조회
	stationID, err := api.GetStationID(stationName)
	if err != nil {
		log.Printf("정류소 ID 조회 실패: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "정류소를 찾을 수 없습니다", "details": err.Error()})
		return
	}

	// 버스 도착 정보 조회
	busArrivals, err := outputs.GetBusInfo(stationID)
	if err != nil {
		log.Printf("버스 정보 조회 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "버스 도착 정보를 가져오는 중 오류가 발생했습니다"})
		return
	}

	// 결과 반환
	c.JSON(http.StatusOK, gin.H{
		"stationName": stationName,
		"busArrivals": busArrivals,
	})
}
