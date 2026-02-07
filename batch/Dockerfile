# 빌드 스테이지
FROM golang:1.24.2 AS builder

WORKDIR /app

# 의존성 파일 복사
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# main.go + 전체 폴더 복사 (config, handlers 등)
RUN go mod tidy
RUN go build -o /batch-service ./cmd/main.go

# 실행 스테이지
#FROM debian:bullseye-slim
FROM debian:bookworm-slim


WORKDIR /root/

# 빌드된 바이너리만 복사
COPY --from=builder /batch-service .


EXPOSE 8081 50001
CMD ["./batch-service"]


# 20251211 Batch
