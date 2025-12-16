package loader

import (
	"context"
)

// 모든 로더가 구현해야 하는 인터페이스
type Loader interface {
	Load(ctx context.Context) error
}

// 여러 로더를 순서대로 실행하는 오케스트레이터
type DataLoader struct {
	loaders []Loader
}

// 가변인자(...)로 여러 Loader 주입 가능
func NewDataLoader(loaders ...Loader) *DataLoader {
	return &DataLoader{loaders: loaders}
}

// 모든 로더 실행
func (d *DataLoader) LoadAllData(ctx context.Context) error {
	for _, loader := range d.loaders {
		if err := loader.Load(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Loader 추가
func (d *DataLoader) Register(loader Loader) {
	d.loaders = append(d.loaders, loader)
}
