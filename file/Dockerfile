# 1. 빌드 스테이지
FROM golang:1.24.2 AS builder

# 빌드 환경 강제 (리눅스용 정적 바이너리)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# 의존성 먼저 복사 (소스 수정 시 캐시 활용)
COPY go.mod go.sum ./
RUN go mod download

# 전체 소스 복사
COPY . .

# [수정] 빌드 대상을 명확히 하고, 불필요한 심볼 제거(-s -w)로 용량 최적화
RUN go build -ldflags '-s -w' -o /file-server ./cmd/main.go

# 2. 실행 스테이지
FROM debian:bookworm-slim

WORKDIR /root/

# 빌드된 바이너리만 복사 (이름을 다르게 해서 찌꺼기와 혼동 방지)
COPY --from=builder /file-server .

# 실행 권한 부여
RUN chmod +x ./file-server

# HTTPS 통신을 위한 인증서 패키지 설치
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


EXPOSE 8091 50001 50002

# CMD를 리스트 형식으로 작성
CMD ["./file-server"]

# 20260205 file