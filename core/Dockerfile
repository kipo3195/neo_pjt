# 빌드 스테이지
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod ./
COPY . .
# main.go + 전체 폴더 복사 (config, handlers 등)
RUN go mod tidy
RUN go build -o main

# 실행 스테이지
#FROM debian:bullseye-slim
FROM debian:bookworm-slim


WORKDIR /root/

# 빌드된 바이너리만 복사
COPY --from=builder /app/main .
COPY ./config ./config 

EXPOSE 8085
CMD ["./main"]


# 20250510 Core
