# 빌드 스테이지
FROM golang:1.24.2 AS builder

WORKDIR /app

# 의존성 파일 복사
COPY go.mod go.sum ./
RUN go mod download

# 프로젝트 전체 소스 복사 
COPY . .

# main.go + 전체 폴더 복사 (config, handlers 등) message-service 바이너리 생성
RUN go build -o /notificator-service ./cmd/main.go

# 실행 스테이지
#FROM debian:bullseye-slim
FROM debian:bookworm-slim

WORKDIR /root/

# 빌드된 바이너리만 복사
COPY --from=builder /notificator-service .

EXPOSE 8082
CMD ["./notificator-service"]


# 20260208 Notificator
