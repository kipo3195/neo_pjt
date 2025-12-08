package job

type Job interface {
	// Execute 메서드는 워커 고루틴에 의해 호출되어 실제 비즈니스 로직을 수행합니다.
	Execute() error
}
