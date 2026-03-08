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

### Problem

기존 Java 기반 WebSocket 서버는 하나의 비동기 스레드가 전체 WebSocket Connection을 순회하며  
Ping 및 메시지를 처리하는 **중앙 집중형(Session Iteration)** 구조였습니다.

이 구조에서는 특정 Connection의 네트워크 지연이 전체 세션 처리에 영향을 줄 수 있으며  
세션 단위의 독립적인 처리와 장애 격리가 어렵다는 문제가 있었습니다.

또한 WebSocket Connection은 **Concurrent Write가 안전하지 않기 때문에**  
클라이언트에게 동시에 메시지 전송 시 Race Condition이 발생할 수 있습니다.

---

### Design

- **Session Isolation**  
  각 WebSocket Connection을 독립적인 Goroutine으로 분리하여 세션 단위 처리 구조로 구성했습니다.  
  [Session Isolation 구현 코드 - 1 connection lifecycle goroutine](https://github.com/kipo3195/neo_pjt/blob/main/notificator/internal/adapter/http/router/notificator_route.go#L40-L47)  
  [Session Isolation 구현 코드 - 2 init Read/Write & ping-pong](https://github.com/kipo3195/neo_pjt/blob/main/notificator/internal/adapter/http/handler/notificator_service_handler.go#L42-L146)  
  
- **Single Writer Pattern**  
  WebSocket의 Concurrent Write 제약을 해결하기 위해 Channel 기반 메시지 전달 구조를 적용했습니다.  
  [Single Writer Pattern 구현 코드 - 채널 데이터 수신](https://github.com/kipo3195/neo_pjt/blob/main/notificator/internal/application/usecase/socket_sender_usecase.go#L55-L83)  

- **Read / Write Goroutine 분리**  
  WebSocket Read는 Blocking I/O이기 때문에 Read / Write 처리를 각각의 Goroutine으로 분리했습니다.  
  [Read / Write Goroutine 분리 구현 코드 - 1 Read](https://github.com/kipo3195/neo_pjt/blob/main/notificator/internal/adapter/http/handler/notificator_service_handler.go#L42-L146)  
  [Read / Write Goroutine 분리 구현 코드 - 2 Write](https://github.com/kipo3195/neo_pjt/blob/main/notificator/internal/application/usecase/socket_sender_usecase.go#L32-L100)  

---

### Data Structure

```go
type SendConnectionEntity struct {
    UserHash string
    Conn     *websocket.Conn
    Chan     chan interface{}
}
```

각 WebSocket 연결은 `SendConnectionEntity`로 관리됩니다.

| Field | Description |
|------|-------------|
| UserHash | 사용자 식별 키 |
| Conn | WebSocket Connection |
| Chan | Write 전용 메시지 채널 |

Producer는 메시지를 Channel로 전달하고,  
각 Session의 Write Goroutine이 이를 WebSocket으로 전송합니다.

---

### Message Flow

```
Producer (Chat / Notification)
   ↓
Session Channel
   ↓
Write Goroutine (per connection)
   ↓
WebSocket Connection
   ↓
Client
```

메시지는 Producer에서 생성된 후 Session Channel을 통해 전달되며  
각 Connection의 Write Goroutine이 WebSocket으로 전송합니다.

---

### Result

- **Concurrent Write 안정성 확보**
  - Channel 기반 Single Writer Pattern으로 Race Condition 제거

- **세션 단위 장애 격리**
  - Connection별 Goroutine 구조로 특정 클라이언트의 네트워크 지연이 전체 시스템에 영향을 주지 않음

- **네트워크 I/O 병목 감소**
  - Read / Write Goroutine 분리를 통해 Blocking I/O 영향 최소화
    
---

## 2. Worker Pool 기반 비동기 이벤트 처리

### Problem
실시간 메시징 환경에서는 **읽지 않은 메시지 수(Unread Count)** 와 같은 이벤트가
짧은 시간 동안 빈번하게 발생할 수 있습니다.

이벤트가 발생할 때마다 즉시 클라이언트로 전송할 경우
불필요한 **Network I/O 증가**와 **이벤트 처리 부하**가 발생할 수 있습니다.

이를 효율적으로 처리하면서도 동시성 문제를 최소화할 수 있는 이벤트 처리 구조가 필요했습니다.

---

### Design
- **Sticky Worker (User Hash 기반)**  
  각 사용자의 이벤트를 특정 Worker에 고정 분배하여, Worker 내부에서는 **Lock 없이 상태 관리** 가능.

- **Timer 기반 Debouncing**  
  이벤트를 일정 시간 동안 모아서 배칭 처리함으로써 클라이언트 전송 횟수를 최소화하고 **I/O 효율 향상**.

- **Worker Pool 구조**  
  Worker Goroutine이 독립적으로 이벤트를 처리하고, 채널(Channel) 기반으로 이벤트를 전달.

---

### Data Structure
- **pendingMap**: 사용자별 상태 관리 (Lock-Free의 핵심)  
- **Event Channel**: Dispatcher → Worker로 이벤트 전달  
- **JobEntity**: 일정 시간 단위로 이벤트 모아서 전송
	
```go
// pendingMap
pendingMap := make(map[string]*entity.ChatCountJobEntity)
```
  
```go
// JobEntity(타이머와 델타값) - 각 워커는 해당 구조체를 value로 갖는 맵 생성
type ChatCountJobEntity struct {
	UserHash string
	Timer    *time.Timer
	Count    int
	Delta    int
}
```
| Field | Description |
|------|-------------|
| UserHash | 룸 키 |
| Timer | 룸 타입 (일반, 오픈) |
| Count | n회 이상 이벤트 수신시 전송 |
| Delta | 미확인 건수 (디바운싱 된 값) |

--- 

### Message Flow

```
미확인 건수 이벤트 발생
   ↓
Hash Dispatcher (UserHash 기반 Worker 선택)
   ↓
Worker Channel 전달
   ↓
Worker Goroutine 내부 처리 (Delta 계산)
   ↓
Timer 기반 (혹은 건수) 배칭
   ↓
Session Channel 
   ↓
Write Goroutine (per connection)
   ↓
WebSocket Connection
   ↓
Client
```

---

### Result

- **Network I/O 최적화**
  - 이벤트를 Timer 기반으로 배칭 처리하여 클라이언트 전송 횟수 감소

- **동시성 처리 효율 향상**
  - User Hash 기반 Worker Sharding을 통해 Lock 없이 이벤트 상태 관리

- **이벤트 처리 안정성 확보**
  - Worker Pool 구조를 통해 이벤트 처리 부하를 분산하고 안정적인 처리 구조 구현
---

## 3. Graceful Shutdown & Resource Lifecycle

### Problem
컨테이너 기반 환경에서 애플리케이션 배포 및 스케일링 시, 서버 종료 시점에 **진행 중인 작업이 중단**되거나 **데이터 손실**이 발생할 수 있습니다.  

안전하게 종료하고 리소스를 정리할 수 있는 **Graceful Shutdown 프로세스**가 필요했습니다.

---

### Design
- **신규 요청 차단**
    - HTTP/gRPC 서버에 더 이상 요청이 들어오지 않도록 막음
- **Worker Pool Job 완료 대기**
    - 처리 중인 이벤트가 안전하게 끝날 때까지 대기
- **외부 연결 종료**
    - Redis, NATS, DB 등 연결을 순차적으로 종료
      
---

### Message Flow
```
SIGTERM / SIGINT 수신
        ↓
HTTP / gRPC 서버 Stop (신규 요청 차단)
        ↓
Worker Pool 진행 중 Job 완료 대기
        ↓
Redis / NATS / DB 연결 종료
```

---

### Result
- **데이터 손실 방지**
    - 진행 중인 이벤트 및 배칭 데이터를 안전하게 처리

- **서비스 안정성 향상**
    - 컨테이너 배포 및 재시작 시 장애 최소화

- **리소스 정리**
    - 외부 연결과 내부 상태를 안전하게 종료하여 메모리/연결 누수 방지

---
## 4. NATS 기반 이벤트 메시징 아키텍처

### Problem
기존 시스템에서는 **Redis Pub/Sub**을 이용하여 서버 노드 간 이벤트 전파를 처리했습니다.

하지만 Redis는 **Single Thread Event Loop 기반 구조**이기 때문에,  
다른 명령 처리나 부하가 높은 작업이 실행될 경우 **Pub/Sub 메시지 전파 지연이 발생할 수 있는 구조적 한계**가 있었습니다.

또한 Redis Pub/Sub는 단순 브로드캐스트 방식으로,  
분산 환경에서 필요한 **Request-Reply 기반의 서비스 간 통신**이나  
**노드 간 상태 동기화 패턴을 구현하기에는 제약이 존재했습니다.**  

이러한 구조를 개선하기 위해 **NATS 기반 메시징 아키텍처**를 도입했습니다.

---

### Design

- **Pub/Sub 패턴**
  - 채팅 메시지 이벤트를 여러 서비스 노드로 실시간 전파

- **Request-Reply 패턴**
  - 채팅방 생성과 같이 **특정 노드에서만 처리해야 하는 요청**이나  
    **노드 상태 조회와 같은 요청에 대한 응답이 필요한 작업**을 처리하기 위해 사용

- **경량 메시지 브로커**
  - 고성능 event routing을 통해 낮은 latency로 메시지 전달

---

### Data Structure
- **Event Payload (JSON / Protobuf)**

```go
type ChatLineEntity struct {
	Cmd           int    `json:"cmd"`
	SendUserHash  string `json:"sendUserHash"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey"`
	Contents      string `json:"contents"`
	SendDate      string `json:"sendDate"`
}
```
| Field | Description |
|------|-------------|
| Cmd | 채팅 이벤트 구분값|
| SendUserHash | 발신자 |
| LineKey | 채팅 메시지의 ID |
| TargetLineKey | 타겟이 되는 채팅 메시지의 ID |
| Contents | 채팅 내용 |
| SendDate | 발신 날짜 |

--- 

### Message Flow

```
Message Service
    ↓
NATS Request (create room event)
    ↓
Notificator Service (Subscriber)
    ↓
NATS Reply         
    ↓   
Message Service Response


Notificator Service
   ↓
NATS Publish (room.created event)
   ↓
All Notificator Services
   ↓
각 노드 Local State 갱신
```

---

### Result

- **실시간 이벤트 전파 안정성 향상**
  - NATS 기반 Pub/Sub을 통해 분산 환경의 Notificator 서비스 간 이벤트를 안정적으로 전파

- **서비스 간 통신 유연성 확보**
  - Request-Reply 패턴을 활용하여 채팅방 생성과 같은 특정 노드 처리 작업 및 상태 조회 요청을 효율적으로 처리

- **분산 환경 상태 동기화**
  - 이벤트 Publish를 통해 모든 Notificator 서비스 노드의 로컬 상태를 일관되게 유지

---

## 5. 분산 데이터 정합성 검증 (File Upload Consistency)

### Problem
채팅 메시지에 파일이 포함되는 경우,  
클라이언트가 파일 업로드를 완료하지 않았거나 업로드 이후 메시지 전송이 실패하면  
스토리지에 **사용되지 않는 파일** 이 남아 **불필요한 스토리지 비용**이 발생할 수 있습니다.

또한 메시지 서비스와 파일 서비스가 분리된 구조에서는  
**파일 업로드 완료 여부와 메시지 전송 간의 데이터 정합성 보장**이 필요했습니다.  

이를 해결하기 위해 **Pre-signed URL 기반 업로드와 Redis 상태 검증을 결합한 정합성 검증 구조**를 설계했습니다.

---

### Design

- **Pre-signed URL 업로드**
  - 파일 서비스에서 업로드 URL과 Transaction ID(TID) 발급

- **Redis 기반 업로드 상태 검증**
  - 파일 업로드 완료 후 TID 기준으로 스토리지 업로드 여부 검증
  - 모든 파일 업로드 완료 시 Redis에 상태 저장

- **Message Service 검증 로직**
  - 채팅 메시지 전송 시 TID를 전달
  - Redis 상태를 확인하여 업로드 완료된 파일만 메시지 처리

- **gRPC 기반 교차 검증**
  - 메시지 전송 시 파일 서비스에 gRPC 호출
  - 해당 파일들의 `sendFlag` 상태 업데이트

- **Batch Cleanup**
  - 업로드 요청만 하고 실제 메시지에 사용되지 않은 파일을 탐지
  - 메시지 서비스와 교차 검증 후 스토리지에서 삭제

---

### Data Structure

```go
type FileUploadUrlHistory struct {
	FileId        string   
	TransactionId string    
	FileName      string    
	ReqUserHash   string    
	UploadUrl     string    
	CreateDate    time.Time 
	UpdateDate    time.Time 
	UploadFlag    string    
	SendFlag      string    
	ErrorFlag     string    
}
```
| Field | Description |
|------|-------------|
| FileId | 파일 구분값 |
| TransactionId | 파일 업로드 요청 구분값 |
| FileName | 파일명 |
| ReqUserHash | 요청 사용자 |
| UploadUrl | 발급된 upload url |
| CreateDate | 파일의 최초 요청일 |
| UpdateDate | 업로드 시간 |
| UploadFlag | 업로드 여부 |
| SendDate | message 서비스에서 전송된 시간 |
| SendFlag | 전송 여부 |
| ErrorFlag | 에러 여부 (스토리지 삭제 대상) |


---
### Message Flow
```
Client
   ↓
File Service (Upload URL + TID 요청)
   ↓
Pre-signed URL 발급

Client
   ↓
Storage 파일 업로드
   ↓
File Service 업로드 완료 요청 (TID)

File Service
   ↓
스토리지 업로드 검증
   ↓
Redis 상태 저장 (Upload Complete)

Client
   ↓
Message Service 채팅 전송 (TID 포함)

Message Service
   ↓
Redis 업로드 상태 검증
   ↓
File Service gRPC 호출
   ↓
sendFlag = Y 업데이트

Batch Service
   ↓
sendFlag = N 파일 탐지
   ↓
Message Service 교차 검증
   ↓
미사용 파일 Storage 삭제
```

---

### Result

- **파일 실재성 보장**
    - Redis 기반 업로드 검증을 통해 실제 업로드된 파일만 메시지 전송에 사용
- **분산 서비스 간 데이터 정합성 확보**
    - Message Service와 File Service 간 gRPC 교차 검증으로 상태 일관성 유지
- **스토리지 비용 최적화**
    - Batch Service의 스케쥴링 작업으로 정합성이 확보되지 않은 파일 정리 및 불필요한 파일 저장 비용 감소

