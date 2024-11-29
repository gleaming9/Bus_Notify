package main

import (
	"github.com/gleaming9/Bus_Notify/api"
	"github.com/gleaming9/Bus_Notify/consume"
	"github.com/gleaming9/Bus_Notify/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	api.LoadStationData()

	// RabbitMQ 서버 실행
	go func() {
		log.Println("RabbitMQ 서버 실행 중...")
		consume.ConsumeFromRabbitMQ()
	}()

	// 라우터 초기화
	router := routes.InitRoutes()

	// 서버 종료 신호 처리 (Graceful Shutdown)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("서버 종료 신호 수신, 종료 중...")
		os.Exit(0)
	}()

	// 서버 실행
	log.Println("서버가 9090 포트에서 실행 중입니다...")
	if err := router.Run(":9090"); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
