# 베이스 이미지
FROM golang:1.23.1 AS builder

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈과 소스 코드 복사
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 go build -o main .

# 실행 이미지
FROM debian:bullseye-slim

WORKDIR /app

# 필수 패키지 설치
RUN apt-get update && apt-get install -y \
    libssl-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get install -y tzdata

# 한국 시간으로 설정
ENV TZ=Asia/Seoul

# 애플리케이션 실행 파일 복사
COPY --from=builder /app/main .

# CSV 파일 복사
COPY bus_stations.csv .

# 포트 노출
EXPOSE 8080

# 실행 명령
CMD ["./main"]
