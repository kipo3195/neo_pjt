package loader

import (
	"context"
	"fmt"
	"log"
)

// Loader: 개별 데이터 로딩 단위
type Loader interface {
	Load(ctx context.Context) error
	Name() string // 로깅이나 디버깅을 위해 로더의 이름을 반환
}

// DataLoader: 여러 로더를 관리하고 실행하는 오케스트레이터
type DataLoader struct {
	loaders []Loader
}

func NewDataLoader(loaders ...Loader) *DataLoader {
	return &DataLoader{loaders: loaders}
}

func (d *DataLoader) Register(l Loader) {
	d.loaders = append(d.loaders, l)
}

func (d *DataLoader) LoadAllData(ctx context.Context) error {
	log.Printf("[DataLoader] Start loading %d items...", len(d.loaders))

	for i, loader := range d.loaders {
		// 인터페이스에 Name() 메서드가 있다면 활용 (없다면 reflect 사용)
		loaderName := fmt.Sprintf("Loader-%d", i)
		if n, ok := any(loader).(interface{ Name() string }); ok {
			loaderName = n.Name()
		}

		log.Printf("[DataLoader] [%d/%d] Running: %s", i+1, len(d.loaders), loaderName)

		// 개별 로더 실행
		if err := loader.Load(ctx); err != nil {
			// 어디서 에러가 났는지 래핑하여 반환
			return fmt.Errorf("[DataLoader] failed at %s: %w", loaderName, err)
		}
	}

	log.Println("[DataLoader] All data loaded successfully.")
	return nil
}
