# Goroutine-based High-Concurrency Real-time Messaging Server

## Overview

Java 기반 메신저 서버의 세션 관리 구조는 동시 연결 수 증가 시 세션 관리 비용과 자원 사용량이 빠르게 증가하는 문제가 있었습니다.  
이를 해결하기 위해 **Go의 Goroutine과 Channel 기반 Concurrency 모델을 활용하여 WebSocket 세션 구조를 재설계한 실시간 메시징 서버**를 구현했습니다.

세션별 독립 Goroutine 구조와 Channel 기반 비동기 메시지 전달, Worker Pool 기반 백그라운드 처리 모델을 통해 **고동시성 환경에서도 안정적인 메시지 처리와 장애 격리를 달성**하는 것을 목표로 설계되었습니다.

---

## Tech Stack

**Language**

- Go

**Networking**

- Gin-Gonic
- Gorilla WebSocket

**Service Communication**

- gRPC
- Protocol Buffers

**Database**

- MariaDB

**Cache**

- Redis

**Message Broker**

- NATS

---

# Technical Challenges

## 1. WebSocket Session Architecture (Concurrent Write Problem)

WebSocket은 **Concurrent Write가 허용되지 않는 Thread-unsafe 구조**이기 때문에, 여러 고루틴에서 동시에 Write를 수행할 경우 Race Condition이 발생할 수 있습니다.

이를 해결하기 위해 **세션별 전용 Write Goroutine과 Channel 기반 메시지 큐 구조**를 설계했습니다.

### Data Structure

```go
type SendConnectionEntity struct {
    UserHash string
    Conn     *websocket.Conn
    Chan     chan interface{}
}
```

### Message Flow 
```
Producer
   ↓
Session Channel
   ↓
Write Goroutine
   ↓
WebSocket Connection
```
### Result
Concurrent Write 문제 해결  
Channel FIFO 기반 메시지 순서 보장  
비즈니스 로직과 네트워크 I/O 분리  
  
## 2. Worker Pool 기반 비동기 이벤트 처리

읽지 않은 메시지 수(Unread Count) 계산과 같은 부하가 큰 로직이 메인 메시지 흐름을 차단하지 않도록 Worker Pool 구조를 설계했습니다.

이벤트는 User Hash 기반으로 특정 Worker에 고정 분배(Sticky Worker) 되도록 설계하여 Worker 내부에서는 Lock 없이 상태를 관리하도록 구현했습니다.

### Event Flow  
```
Event
   ↓
Hash Dispatcher
   ↓
Worker Channel
   ↓
Worker Goroutine
```
  
### Core Strategy
- User Hash 기반 Worker Sharding  
- Worker Local Map 기반 Lock-Free 상태 관리  
- Timer 기반 Debouncing으로 DB 업데이트 배치 처리  

## 3. Full-Duplex WebSocket Architecture

WebSocket 연결은 Read와 Write를 분리된 Goroutine으로 운영하는 Full-Duplex 구조로 설계했습니다.  
```
Connection Structure
WebSocket Connection
├── Read Goroutine
└── Write Goroutine
```  
### Core Strategy 
- 네트워크 지연으로 Write가 블로킹되더라도 Read 처리나 다른 세션에 영향을 주지 않도록 장애 격리를 구현  
- Ping/Pong 기반 Heartbeat를 구현하여 응답이 없는 좀비 커넥션을 자동 정리하도록 처리  

## 4. Graceful Shutdown & Resource Lifecycle

컨테이너 기반 환경에서 안전한 배포를 위해 Graceful Shutdown 프로세스를 구현했습니다.

Shutdown Flow
SIGTERM / SIGINT
     ↓
HTTP / gRPC Server Stop (신규 요청 차단)
     ↓
Worker Pool Job 완료 대기
     ↓
In-Memory Buffer Flush
     ↓
Redis / NATS / DB 연결 종료

### Core Strategy 
- 이를 통해 서버 종료 시 메시지 데이터 유실을 방지하도록 설계  
## 5. Distributed Data Consistency (File & Message Service)

파일 업로드와 메시지 전송이 서로 다른 서비스에서 처리되는 분산 환경에서 데이터 정합성을 보장하기 위해 다음 구조를 설계했습니다.
```
Data Consistency Flow
Client
   ↓
Pre-signed URL File Upload
   ↓
Message Service (Transaction ID 전달)
   ↓
gRPC Validation
   ↓
File Service 상태 검증
```
### Core Strategy 
- gRPC를 활용하여 메시지 전송 완료 후 File 서비스의 상태를 즉각 업데이트하여 최종 일관성 확보
- Batch 검증 프로세스를 통해 전송되지 않은 잔여 파일을 자동 정리하여 스토리지 비용을 최적화
