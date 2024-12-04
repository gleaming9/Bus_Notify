package service

import (
	"fmt"
	"github.com/gleaming9/Bus_Notify/api"
	"github.com/gleaming9/Bus_Notify/model"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func MonitorBusArrival(req model.AlertRequest) {
	// 정류소 ID와 버스 도착 정보 조회
	stationID, err := api.GetStationID(req.StationName)
	if err != nil {
		log.Printf("정류소 ID 조회 실패: %v", err)
		return
	}

	for {
		// 버스 도착 시간 확인 및 메일 발송 조건 체크
		arrivalInfo, err := api.GetBusArrivalInfo(stationID)
		if err != nil {
			log.Printf("버스 도착 정보 조회 실패: %v", err)
			return
		}

		busData := [20]model.BusData{}
		cnt := 0
		for _, bus := range arrivalInfo.Body.BusArrivalList {
			routeInfo, err := api.GetBusRouteInfo(bus.RouteID) // 노선 ID를 기반으로 버스 노선 정보 가져오기
			if err != nil {
				log.Printf("노선 ID %s의 추가 정보를 가져오는 데 실패했습니다: %v\n", bus.RouteID, err)
				continue
			}

			if bus.PredictTime1 != "" || bus.PredictTime1 < "30" {
				busData[cnt] = model.BusData{
					First:  routeInfo.MsgBody.BusRouteInfoItem.RouteName,
					Second: bus.PredictTime1,
				}
				cnt++
			} else if bus.PredictTime2 != "" || bus.PredictTime1 < "30" {
				busData[cnt] = model.BusData{
					First:  routeInfo.MsgBody.BusRouteInfoItem.RouteName,
					Second: bus.PredictTime2,
				}
				cnt++
			}
		}

		routeName, timeleft := checkArrivalCondition(busData[:cnt], req.TargetTime)
		if routeName == "" && timeleft == 0 {
			// 이메일 발송

			subject := fmt.Sprintf("%s 버스 도착 실패 알림", req.StationName)
			body := fmt.Sprintf("설정하신 시간 이내에 도착하는 버스가 없습니다.")
			if err := sendEmail(req.Email, subject, body); err != nil {
				log.Printf("이메일 전송 실패: %v", err)
				return
			}
			log.Printf("이메일 전송 완료: %s", req.Email)

			return
		} else if routeName == "" && timeleft == 1 {
			log.Printf("도착 예정 버스가 없습니다.")
			time.Sleep(1 * time.Minute)
			continue
		} else {
			// 이메일 발송
			// 문자열을 파싱하고 시간 계산
			parts := strings.Split(req.TargetTime, ":")
			h, _ := strconv.Atoi(parts[0])
			m, _ := strconv.Atoi(parts[1])
			newTime := time.Date(0, 1, 1, h, m, 0, 0, time.Local).Add(-time.Duration(timeleft) * time.Minute)

			subject := fmt.Sprintf("%s 버스 도착 예정 알림", req.StationName)
			body := fmt.Sprintf("도착 예정 버스: %s\n %s 까지 도착해야합니다!", routeName, newTime.Format("15:04"))
			if err := sendEmail(req.Email, subject, body); err != nil {
				log.Printf("이메일 전송 실패: %v", err)
				return
			}
			log.Printf("이메일 전송 완료: %s", req.Email)
			return
		}
	}
}

// checkArrivalCondition: 특정 조건을 만족하는지 확인
func checkArrivalCondition(busData []model.BusData, targetTime string) (string, int) {
	now := time.Now()
	// 현재 시간과 분 추출
	hours := now.Hour()
	minutes := now.Minute()
	nowTime_minutes := (hours * 60) + minutes

	parts := strings.Split(targetTime, ":")
	// 시간과 분을 정수로 변환
	hours, err1 := strconv.Atoi(parts[0])
	minutes, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		log.Println("시간 또는 분을 숫자로 변환할 수 없습니다.")
		return "", 0
	}
	// 총 분 계산
	targetTime_minutes := (hours * 60) + minutes

	log.Printf("타겟 시간 : %d\n", targetTime_minutes)

	closestTime := 0
	closestBus := ""
	for _, bus := range busData {
		b, err := strconv.Atoi(bus.Second)
		if err != nil {
			log.Println("도착 예상 시간을 숫자로 변환할 수 없습니다")
			return "", 0
		}
		busTime_minutes := nowTime_minutes + b
		log.Printf("%s 버스 시간 : %d\n", bus.First, busTime_minutes)
		// targetTime_minutes보다 작은 값 중 가장 근접한 값 찾기
		if (busTime_minutes < targetTime_minutes) && (targetTime_minutes-busTime_minutes < targetTime_minutes-closestTime) {
			closestTime = busTime_minutes
			closestBus = bus.First
		}
	}
	log.Printf("가장 가까운 버스 %s 시간 : %d\n", closestBus, closestTime)
	if closestTime == 0 {
		return "", 0
		// 가장 근접한 도착 시간이 15분 이내면 알림
	} else if targetTime_minutes-closestTime <= 15 {
		return closestBus, targetTime_minutes - closestTime
	} else {
		return "", 1
	}
}

// 이메일 전송 함수
func sendEmail(to, subject, body string) error {
	from := "mathasdf0@gmail.com"
	password := "vxrq llox tohy brca"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 이메일 메시지 작성
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// SMTP 연결 및 이메일 전송
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Fatalf("이메일 전송 실패: %v", err)
		return err
	}

	fmt.Println("이메일 전송 성공!")
	return nil
}
